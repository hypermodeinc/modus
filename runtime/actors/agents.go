/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package actors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/messages"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
	"github.com/rs/xid"

	wasm "github.com/tetratelabs/wazero/api"
	goakt "github.com/tochemey/goakt/v3/actor"
)

type AgentInfo struct {
	Id     string
	Name   string
	Status AgentStatus
}

type AgentStatus = string

// TODO: validate these statuses are needed and used correctly
const (
	AgentStatusUninitialized AgentStatus = "uninitialized"
	AgentStatusError         AgentStatus = "error"
	AgentStatusStarting      AgentStatus = "starting"
	AgentStatusStarted       AgentStatus = "started"
	AgentStatusStopping      AgentStatus = "stopping"
	AgentStatusStopped       AgentStatus = "stopped"
	AgentStatusSuspended     AgentStatus = "suspended"
	AgentStatusTerminated    AgentStatus = "terminated"
)

func SpawnAgentActor(ctx context.Context, agentName string) (*AgentInfo, error) {
	plugin, ok := plugins.GetPluginFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no plugin found in context")
	}

	agentId := xid.New().String()
	host := wasmhost.GetWasmHost(ctx)

	// We spawn the actor in a goroutine to avoid blocking the host function while the actor is being spawned.
	go func() {
		// The actor needs to spawn in a new context. Otherwise, when the original context is cancelled (such as when the function completes),
		// the actor initialization would be cancelled too.
		ctx := context.Background()
		ctx = context.WithValue(ctx, utils.WasmHostContextKey, host)
		ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)
		ctx = context.WithValue(ctx, utils.AgentIdContextKey, agentId)
		ctx = context.WithValue(ctx, utils.AgentNameContextKey, agentName)

		actor := NewWasmAgentActor(agentId, agentName, plugin)
		actorName := fmt.Sprintf("agent-%s", agentId)
		if _, err := _actorSystem.Spawn(ctx, actorName, actor); err != nil {
			logger.Err(ctx, err).Msg("Error spawning actor.")
		}
	}()

	info := &AgentInfo{
		Id:     agentId,
		Name:   agentName,
		Status: AgentStatusStarting,
	}

	return info, nil
}

type AgentMessageResponse struct {
	Data *string
}

func SendAgentMessage(ctx context.Context, agentId string, msgName string, data *string, timeout int64) (*AgentMessageResponse, error) {

	addr, pid, err := _actorSystem.ActorOf(ctx, getActorName(agentId))
	if err != nil {
		if strings.HasSuffix(err.Error(), " not found") {
			return nil, fmt.Errorf("agent %s not found", agentId)
		} else {
			return nil, fmt.Errorf("error getting actor for agent %s: %w", agentId, err)
		}
	}

	_ = addr // TODO: this will be used when we implement remote actors with clustering

	msg := &messages.AgentRequestMessage{
		Name:    msgName,
		Data:    data,
		Respond: timeout > 0,
	}

	if timeout == 0 {
		if err := goakt.Tell(ctx, pid, msg); err != nil {
			return nil, fmt.Errorf("error sending message to agent %s: %w", pid.ID(), err)
		}
		return &AgentMessageResponse{}, nil
	}

	res, err := goakt.Ask(ctx, pid, msg, time.Duration(timeout))
	if err != nil {
		return nil, fmt.Errorf("error sending message to agent %s: %w", pid.ID(), err)
	}

	if response, ok := res.(*messages.AgentResponseMessage); ok {
		return &AgentMessageResponse{response.Data}, nil
	} else {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

func getActorName(agentId string) string {
	return "agent-" + agentId
}

type WasmAgentActor struct {
	agentId   string
	agentName string
	plugin    *plugins.Plugin
	host      wasmhost.WasmHost
	module    wasm.Module
	buffers   utils.OutputBuffers
}

func NewWasmAgentActor(agentId, agentName string, plugin *plugins.Plugin) *WasmAgentActor {
	return &WasmAgentActor{
		agentId:   agentId,
		agentName: agentName,
		plugin:    plugin,
	}
}

func (a *WasmAgentActor) PreStart(ac *goakt.Context) error {
	ctx := ac.Context()

	logger.Info(ctx).Bool("user_visible", true).Msg("Starting agent.")
	start := time.Now()

	a.host = wasmhost.GetWasmHost(ctx)
	a.buffers = utils.NewOutputBuffers()
	if mod, err := a.host.GetModuleInstance(ctx, a.plugin, a.buffers); err != nil {
		return err
	} else {
		a.module = mod
	}

	if err := a.activateAgent(ctx, false); err != nil {
		logger.Err(ctx, err).Bool("user_visible", true).Msg("Error starting agent.")
		return err
	}

	duration := time.Since(start)
	logger.Info(ctx).Dur("duration_ms", duration).Bool("user_visible", true).Msg("Agent started successfully.")
	return nil
}

func (a *WasmAgentActor) PostStop(ac *goakt.Context) error {
	ctx := ac.Context()
	defer a.module.Close(ctx)

	// the context may not have these values set
	ctx = context.WithValue(ctx, utils.PluginContextKey, a.plugin)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, a.agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, a.agentName)

	logger.Info(ctx).Bool("user_visible", true).Msg("Stopping agent.")
	start := time.Now()

	if err := a.shutdownAgent(ctx); err != nil {
		logger.Err(ctx, err).Bool("user_visible", true).Msg("Error stopping agent.")
		return err
	}

	duration := time.Since(start)
	logger.Info(ctx).Dur("duration_ms", duration).Bool("user_visible", true).Msg("Agent stopped successfully.")
	return nil
}

func (a *WasmAgentActor) Receive(rc *goakt.ReceiveContext) {
	// Create a new context to avoid having the actor handle the message in the host function context.
	// This is important for async messages, where the calling function may return before the message is handled.
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.WasmHostContextKey, a.host)
	ctx = context.WithValue(ctx, utils.PluginContextKey, a.plugin)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, a.agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, a.agentName)

	switch msg := rc.Message().(type) {
	case *messages.AgentRequestMessage:

		logger.Info(ctx).Str("msg_name", msg.Name).Msg("Received message.")
		start := time.Now()

		fnInfo, err := a.host.GetFunctionInfo("_modus_agent_handle_message")
		if err != nil {
			rc.Err(err)
			return
		}

		params := map[string]any{
			"msgName": msg.Name,
			"data":    msg.Data,
		}

		execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
		if err != nil {
			rc.Err(err)
			return
		}

		if msg.Respond {
			result := execInfo.Result()
			response := &messages.AgentResponseMessage{}
			if result != nil {
				switch result := result.(type) {
				case string:
					response.Data = &result
				case *string:
					response.Data = result
				default:
					err := fmt.Errorf("unexpected result type: %T", result)
					logger.Err(ctx, err).Msg("Error handling message.")
					rc.Err(err)
					return
				}
			}
			rc.Response(response)
		}

		duration := time.Since(start)
		logger.Info(ctx).Str("msg_name", msg.Name).Dur("duration_ms", duration).Msg("Message handled successfully.")

	default:
		rc.Unhandled()
	}
}

func (a *WasmAgentActor) activateAgent(ctx context.Context, reloading bool) error {

	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_activate")
	if err != nil {
		return err
	}

	params := map[string]any{
		"name":      a.agentName,
		"id":        a.agentId,
		"reloading": reloading,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)

	_ = execInfo // TODO
	return err
}

func (a *WasmAgentActor) shutdownAgent(ctx context.Context) error {

	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_shutdown")
	if err != nil {
		return err
	}
	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, nil)

	_ = execInfo // TODO
	return err
}

func (a *WasmAgentActor) getAgentState(ctx context.Context) (*string, error) {

	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_get_state")
	if err != nil {
		return nil, err
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, nil)
	if err != nil {
		return nil, err
	}

	result := execInfo.Result()
	if result == nil {
		return nil, nil
	}

	switch state := result.(type) {
	case string:
		return &state, nil
	case *string:
		return state, nil
	default:
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
}

func (a *WasmAgentActor) setAgentState(ctx context.Context, data *string) error {
	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_set_state")
	if err != nil {
		return err
	}

	params := map[string]any{
		"data": data,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
	if err != nil {
		return err
	}

	_ = execInfo // TODO
	return err
}

func (a *WasmAgentActor) reloadModule(ctx context.Context, plugin *plugins.Plugin) error {

	// the context may not have these values set
	ctx = context.WithValue(ctx, utils.PluginContextKey, a.plugin)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, a.agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, a.agentName)

	logger.Info(ctx).Bool("user_visible", true).Msg("Reloading module for agent.")

	// get the current state and close the module instance
	state, err := a.getAgentState(ctx)
	if err != nil {
		logger.Err(ctx, err).Bool("user_visible", true).Msg("Error getting agent state.")
		return err
	}
	a.module.Close(ctx)

	// create a new module instance and assign it to the actor
	a.plugin = plugin
	a.buffers = utils.NewOutputBuffers()
	mod, err := a.host.GetModuleInstance(ctx, a.plugin, a.buffers)
	if err != nil {
		return err
	}
	a.module = mod

	// activate the agent in the new module instance
	if err := a.activateAgent(ctx, true); err != nil {
		logger.Err(ctx, err).Bool("user_visible", true).Msg("Error reloading agent.")
		return err
	}

	// restore the state in the new module instance
	if err := a.setAgentState(ctx, state); err != nil {
		logger.Err(ctx, err).Bool("user_visible", true).Msg("Error setting agent state.")
		return err
	}

	logger.Info(ctx).Bool("user_visible", true).Msg("Agent reloaded module successfully.")

	return nil
}
