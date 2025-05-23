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

	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/messages"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
	"github.com/rs/xid"

	wasm "github.com/tetratelabs/wazero/api"
	goakt "github.com/tochemey/goakt/v3/actor"
)

type agentInfo struct {
	id     string
	name   string
	status agentStatus
}

type agentStatus = string

const (
	agentStatusStarting   agentStatus = "starting"
	agentStatusRunning    agentStatus = "running"
	agentStatusSuspending agentStatus = "suspending"
	agentStatusSuspended  agentStatus = "suspended"
	agentStatusResuming   agentStatus = "resuming"
	agentStatusStopping   agentStatus = "stopping"
	agentStatusTerminated agentStatus = "terminated"
)

func StartAgent(ctx context.Context, agentName string) (*agentInfo, error) {
	plugin, ok := plugins.GetPluginFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no plugin found in context")
	}

	agentId := xid.New().String()
	host := wasmhost.GetWasmHost(ctx)
	spawnActorForAgent(host, plugin, agentId, agentName, false, nil)

	info := &agentInfo{
		id:     agentId,
		name:   agentName,
		status: agentStatusStarting,
	}

	return info, nil
}

func spawnActorForAgent(host wasmhost.WasmHost, plugin *plugins.Plugin, agentId, agentName string, resuming bool, initialState *string) {
	// The actor needs to spawn in its own context, so we don't pass one in to this function.
	// If we did, then when the original context was cancelled or completed, the actor initialization would be cancelled too.

	// We spawn the actor in a goroutine to avoid blocking while the actor is being spawned.
	// This allows many agents to be spawned in parallel, if needed.
	go func() {
		ctx := context.Background()
		ctx = context.WithValue(ctx, utils.WasmHostContextKey, host)
		ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)
		ctx = context.WithValue(ctx, utils.AgentIdContextKey, agentId)
		ctx = context.WithValue(ctx, utils.AgentNameContextKey, agentName)

		actor := newWasmAgentActor(agentId, agentName, host, plugin)
		actorName := getActorName(agentId)

		if resuming {
			actor.status = agentStatusResuming
			actor.initialState = initialState
		} else {
			actor.status = agentStatusStarting
		}

		if _, err := _actorSystem.Spawn(ctx, actorName, actor); err != nil {
			logger.Err(ctx, err).Msg("Error spawning actor for agent.")
		}
	}()
}

func StopAgent(ctx context.Context, agentId string) bool {
	pid, err := getActorPid(ctx, agentId)
	if err != nil {
		logger.Err(ctx, err).Msg("Error stopping agent.")
		return false
	}

	actor := pid.Actor().(*wasmAgentActor)
	actor.status = agentStatusStopping

	if err := pid.Shutdown(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error stopping agent.")
		return false
	}

	return true
}

type agentMessageResponse struct {
	Data *string
}

func SendAgentMessage(ctx context.Context, agentId string, msgName string, data *string, timeout int64) (*agentMessageResponse, error) {

	pid, err := getActorPid(ctx, agentId)
	if err != nil {
		return nil, err
	}

	msg := &messages.AgentRequestMessage{
		Name:    msgName,
		Data:    data,
		Respond: timeout > 0,
	}

	if timeout == 0 {
		if err := goakt.Tell(ctx, pid, msg); err != nil {
			return nil, fmt.Errorf("error sending message to agent %s: %w", pid.ID(), err)
		}
		return &agentMessageResponse{}, nil
	}

	res, err := goakt.Ask(ctx, pid, msg, time.Duration(timeout))
	if err != nil {
		return nil, fmt.Errorf("error sending message to agent %s: %w", pid.ID(), err)
	}

	if response, ok := res.(*messages.AgentResponseMessage); ok {
		return &agentMessageResponse{response.Data}, nil
	} else {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

func getActorName(agentId string) string {
	return "agent-" + agentId
}

func getActorPid(ctx context.Context, agentId string) (*goakt.PID, error) {

	addr, pid, err := _actorSystem.ActorOf(ctx, getActorName(agentId))
	if err != nil {
		if strings.HasSuffix(err.Error(), " not found") {
			return nil, fmt.Errorf("agent %s not found", agentId)
		} else {
			return nil, fmt.Errorf("error getting actor for agent %s: %w", agentId, err)
		}
	}

	_ = addr // TODO: this will be used when we implement remote actors with clustering

	return pid, nil
}

type wasmAgentActor struct {
	agentId      string
	agentName    string
	status       agentStatus
	plugin       *plugins.Plugin
	host         wasmhost.WasmHost
	module       wasm.Module
	buffers      utils.OutputBuffers
	initialState *string
}

func newWasmAgentActor(agentId, agentName string, host wasmhost.WasmHost, plugin *plugins.Plugin) *wasmAgentActor {
	return &wasmAgentActor{
		agentId:   agentId,
		agentName: agentName,
		host:      host,
		plugin:    plugin,
	}
}

func (a *wasmAgentActor) PreStart(ac *goakt.Context) error {
	ctx := a.newContext()

	switch a.status {
	case agentStatusStarting:
		logger.Info(ctx).Msg("Starting agent.")
	case agentStatusResuming, agentStatusSuspended:
		a.status = agentStatusResuming
		logger.Info(ctx).Msg("Resuming agent.")
	default:
		return fmt.Errorf("invalid agent status for actor PreStart: %s", a.status)
	}

	if err := a.saveState(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error saving agent state.")
	}

	start := time.Now()

	a.buffers = utils.NewOutputBuffers()
	if mod, err := a.host.GetModuleInstance(ctx, a.plugin, a.buffers); err != nil {
		return err
	} else {
		a.module = mod
	}

	if err := a.activateAgent(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error activating agent.")
		return err
	}

	if a.status == agentStatusResuming {
		if err := a.setAgentState(ctx, a.initialState); err != nil {
			logger.Err(ctx, err).Msg("Error resuming agent state.")
		}
		a.initialState = nil
	}

	duration := time.Since(start)
	if a.status == agentStatusResuming {
		logger.Info(ctx).Msg("Agent resumed successfully.")
	} else {
		logger.Info(ctx).Dur("duration_ms", duration).Msg("Agent started successfully.")
	}

	a.status = agentStatusRunning

	if err := a.saveState(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error saving agent state.")
	}

	return nil
}

func (a *wasmAgentActor) PostStop(ac *goakt.Context) error {
	ctx := a.newContext()
	defer a.module.Close(ctx)

	switch a.status {
	case agentStatusRunning, agentStatusSuspending:
		a.status = agentStatusSuspending
		logger.Info(ctx).Msg("Suspending agent.")
	case agentStatusStopping:
		logger.Info(ctx).Msg("Stopping agent.")

	default:
		return fmt.Errorf("invalid agent status for actor PostStop: %s", a.status)
	}

	if err := a.saveState(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error saving agent state.")
	}

	start := time.Now()

	if err := a.shutdownAgent(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error shutting down agent.")
		return err
	}

	duration := time.Since(start)
	switch a.status {
	case agentStatusSuspending:
		a.status = agentStatusSuspended
		if err := a.saveState(ctx); err != nil {
			return err
		}
		logger.Info(ctx).Msg("Agent suspended successfully.")
	case agentStatusStopping:
		a.status = agentStatusTerminated
		if err := a.saveState(ctx); err != nil {
			return err
		}
		logger.Info(ctx).Dur("duration_ms", duration).Msg("Agent terminated successfully.")
	default:
		return fmt.Errorf("invalid agent status for actor PostStop: %s", a.status)
	}

	return nil
}

func (a *wasmAgentActor) saveState(ctx context.Context) error {
	var data string
	if a.module != nil {
		if d, err := a.getAgentState(ctx); err != nil {
			return fmt.Errorf("error getting state from agent: %w", err)
		} else if d != nil {
			data = *d
		}
	}

	if err := db.WriteAgentState(ctx, db.AgentState{
		Id:        a.agentId,
		Name:      a.agentName,
		Status:    a.status,
		Data:      data,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("error saving state to database: %w", err)
	}

	return nil
}

func (a *wasmAgentActor) Receive(rc *goakt.ReceiveContext) {
	ctx := a.newContext()

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

		// save the state after handling the message to ensure the state is up to date in case of hard termination
		if err := a.saveState(ctx); err != nil {
			logger.Err(ctx, err).Msg("Error saving agent state.")
		}

	default:
		rc.Unhandled()
	}
}

func (a *wasmAgentActor) newContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.WasmHostContextKey, a.host)
	ctx = context.WithValue(ctx, utils.PluginContextKey, a.plugin)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, a.agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, a.agentName)
	return ctx
}

func (a *wasmAgentActor) activateAgent(ctx context.Context) error {

	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_activate")
	if err != nil {
		return err
	}

	params := map[string]any{
		"name":      a.agentName,
		"id":        a.agentId,
		"reloading": a.status == agentStatusResuming,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)

	_ = execInfo // TODO
	return err
}

func (a *wasmAgentActor) shutdownAgent(ctx context.Context) error {

	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_shutdown")
	if err != nil {
		return err
	}

	params := map[string]any{
		"suspending": a.status == agentStatusSuspending,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)

	_ = execInfo // TODO
	return err
}

func (a *wasmAgentActor) getAgentState(ctx context.Context) (*string, error) {

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

func (a *wasmAgentActor) setAgentState(ctx context.Context, data *string) error {
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

func (a *wasmAgentActor) reloadModule(ctx context.Context, plugin *plugins.Plugin) error {

	// the context may not have these values set
	ctx = context.WithValue(ctx, utils.PluginContextKey, a.plugin)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, a.agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, a.agentName)

	logger.Info(ctx).Msg("Reloading module for agent.")

	a.status = agentStatusSuspending
	if err := a.shutdownAgent(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error shutting down agent.")
		return err
	}

	// get the current state and close the module instance
	state, err := a.getAgentState(ctx)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting agent state.")
		return err
	}
	a.module.Close(ctx)
	a.status = agentStatusSuspended

	// create a new module instance and assign it to the actor
	a.plugin = plugin
	a.buffers = utils.NewOutputBuffers()
	mod, err := a.host.GetModuleInstance(ctx, a.plugin, a.buffers)
	if err != nil {
		return err
	}
	a.module = mod

	// activate the agent in the new module instance
	a.status = agentStatusResuming
	if err := a.activateAgent(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error reloading agent.")
		return err
	}

	// resume the state in the new module instance
	if err := a.setAgentState(ctx, state); err != nil {
		logger.Err(ctx, err).Msg("Error setting agent state.")
		return err
	}

	a.status = agentStatusRunning
	logger.Info(ctx).Msg("Agent reloaded module successfully.")

	return nil
}
