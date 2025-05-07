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
	"time"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/messages"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
	"github.com/puzpuzpuz/xsync/v4"

	wasm "github.com/tetratelabs/wazero/api"
	goakt "github.com/tochemey/goakt/v3/actor"
)

func Activate(ctx context.Context, plugin *plugins.Plugin) error {

	host := wasmhost.GetWasmHost(ctx)
	fnInfo, _ := host.GetFunctionInfo("_modus_register_agents")
	if fnInfo == nil {
		// No agents to register
		return nil
	}

	buffers := utils.NewOutputBuffers()

	mod, err := host.GetModuleInstance(ctx, plugin, buffers)
	if err != nil {
		return fmt.Errorf("error getting module instance: %w", err)
	}
	defer mod.Close(ctx)

	if _, err := host.CallFunctionInModule(ctx, mod, buffers, fnInfo, nil); err != nil {
		return fmt.Errorf("error registering agents: %w", err)
	}

	agentRegistry.Range(func(key string, agent *Agent) bool {
		if agent.Pid == nil {
			err = agent.spawnActor(ctx, plugin)
		} else {
			err = agent.reloadModule(ctx, plugin)
		}
		return err != nil
	})
	if err != nil {
		return err
	}

	return nil
}

type Agent struct {
	Id   int32
	Pid  *goakt.PID
	Name string
}

func (agent *Agent) spawnActor(ctx context.Context, plugin *plugins.Plugin) error {
	host := wasmhost.GetWasmHost(ctx)
	buffers := utils.NewOutputBuffers()
	mod, err := host.GetModuleInstance(ctx, plugin, buffers)
	if err != nil {
		return err
	}

	actor := NewWasmAgentActor(agent, host, mod, buffers)
	pid, err := _actorSystem.Spawn(ctx, agent.Name, actor)
	if err != nil {
		return fmt.Errorf("error spawning actor for '%s' agent: %w", agent.Name, err)
	}

	agent.Pid = pid
	return nil
}

func (agent *Agent) reloadModule(ctx context.Context, plugin *plugins.Plugin) error {
	actor := agent.Pid.Actor().(*WasmAgentActor)

	logger.Info(ctx).Str("agent", agent.Name).Int32("agentId", agent.Id).Bool("user_visible", true).Msg("Reloading module for agent.")

	// get the current state and close the module
	state, err := actor.getAgentState(ctx)
	if err != nil {
		return fmt.Errorf("error getting agent state: %w", err)
	}
	actor.module.Close(ctx)

	// create a new module instance and assign it to the actor
	host := wasmhost.GetWasmHost(ctx)
	buffers := utils.NewOutputBuffers()
	mod, err := host.GetModuleInstance(ctx, plugin, buffers)
	if err != nil {
		return err
	}
	actor.module = mod

	// restore the state in the new module
	if err := actor.setAgentState(ctx, *state); err != nil {
		return fmt.Errorf("error setting agent state: %w", err)
	}

	return nil
}

var agentRegistry = xsync.NewMap[string, *Agent]()

func getAgent(ctx context.Context, id int32) (*Agent, error) {
	if key, err := getAgentKey(ctx, id); err != nil {
		return nil, err
	} else if agent, ok := agentRegistry.Load(key); ok {
		return agent, nil
	} else {
		return nil, fmt.Errorf("agent with id %d not found", id)
	}
}

func getActorForAgent(ctx context.Context, agentId int32) (*WasmAgentActor, error) {
	agent, err := getAgent(ctx, agentId)
	if err != nil {
		return nil, err
	}
	actor := agent.Pid.Actor()
	if actor == nil {
		return nil, fmt.Errorf("actor for agent %d not found", agentId)
	}
	wasmActor, ok := actor.(*WasmAgentActor)
	if !ok {
		return nil, fmt.Errorf("actor for agent %d is not a WasmAgentActor", agentId)
	}
	return wasmActor, nil
}

func getAgentKey(ctx context.Context, agentId int32) (string, error) {
	if plugin, ok := plugins.GetPluginFromContext(ctx); !ok {
		return "", fmt.Errorf("no plugin found in context")
	} else {
		return fmt.Sprintf("%s:%d", plugin.Name(), agentId), nil
	}
}

func RegisterAgent(ctx context.Context, agentId int32, name string) error {
	key, err := getAgentKey(ctx, agentId)
	if err != nil {
		return err
	}

	agentRegistry.LoadOrStore(key, &Agent{
		Id:   agentId,
		Name: name,
	})

	// actual, found := agentRegistry.LoadAndStore(key, &Agent{
	// 	Id:   agentId,
	// 	Name: name,
	// })

	// // If the actor already exists, we need to stop it before spawning a new one.
	// if found {
	// 	if err := actual.Pid.Shutdown(ctx); err != nil {
	// 		return fmt.Errorf("error shutting down existing agent %d: %w", agentId, err)
	// 	}
	// }

	return nil
}

func SendAgentMessage(ctx context.Context, agentId int32, msgName string, data *string, timeout int64) (*string, error) {
	agent, err := getAgent(ctx, agentId)
	if err != nil {
		return nil, err
	}
	if agent.Pid == nil {
		return nil, fmt.Errorf("actor for agent %d not found", agentId)
	}

	msg := &messages.AgentRequestMessage{
		Name:    msgName,
		Data:    data,
		Respond: timeout > 0,
	}

	if timeout == 0 {
		if err := goakt.Tell(ctx, agent.Pid, msg); err != nil {
			return nil, fmt.Errorf("error sending message to agent %s: %w", agent.Pid.ID(), err)
		}
		return nil, nil
	}

	res, err := goakt.Ask(ctx, agent.Pid, msg, time.Duration(timeout))
	if err != nil {
		return nil, fmt.Errorf("error sending message to agent %s: %w", agent.Pid.ID(), err)
	}

	if response, ok := res.(*messages.AgentResponseMessage); ok {
		return response.Data, nil
	} else {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

func GetAgentState(ctx context.Context, agentId int32) (*string, error) {
	actor, err := getActorForAgent(ctx, agentId)
	if err != nil {
		return nil, err
	}
	state, err := actor.getAgentState(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting agent state: %w", err)
	}
	return state, nil
}

func SetAgentState(ctx context.Context, agentId int32, data string) error {
	actor, err := getActorForAgent(ctx, agentId)
	if err != nil {
		return err
	}

	if err := actor.setAgentState(ctx, data); err != nil {
		return fmt.Errorf("error setting agent state: %w", err)
	}
	return nil
}

type WasmAgentActor struct {
	agent   *Agent
	host    wasmhost.WasmHost
	module  wasm.Module
	buffers utils.OutputBuffers
}

func NewWasmAgentActor(agent *Agent, host wasmhost.WasmHost, mod wasm.Module, buffers utils.OutputBuffers) *WasmAgentActor {
	return &WasmAgentActor{agent, host, mod, buffers}
}

func (a *WasmAgentActor) PreStart(ac *goakt.Context) error {
	ctx := ac.Context()

	logger.Info(ctx).Str("agent", a.agent.Name).Int32("agentId", a.agent.Id).Bool("user_visible", true).Msg("Starting agent...")
	if err := a.startAgent(ctx); err != nil {
		logger.Err(ctx, err).Str("agent", a.agent.Name).Int32("agentId", a.agent.Id).Bool("user_visible", true).Msg("Error starting agent.")
		return err
	}

	logger.Info(ctx).Str("agent", a.agent.Name).Int32("agentId", a.agent.Id).Bool("user_visible", true).Msg("Agent started successfully.")
	return nil
}

func (a *WasmAgentActor) PostStop(ac *goakt.Context) error {
	ctx := ac.Context()
	defer a.module.Close(ctx)

	logger.Info(ctx).Str("agent", a.agent.Name).Int32("agentId", a.agent.Id).Bool("user_visible", true).Msg("Stopping agent...")
	if err := a.stopAgent(ctx); err != nil {
		logger.Err(ctx, err).Str("agent", a.agent.Name).Int32("agentId", a.agent.Id).Bool("user_visible", true).Msg("Error stopping agent.")
		return err
	}

	logger.Info(ctx).Str("agent", a.agent.Name).Int32("agentId", a.agent.Id).Bool("user_visible", true).Msg("Agent stopped successfully.")
	return nil
}

func (a *WasmAgentActor) Receive(rc *goakt.ReceiveContext) {
	ctx := rc.Context()

	switch msg := rc.Message().(type) {
	case *messages.AgentRequestMessage:

		fnInfo, err := a.host.GetFunctionInfo("_modus_agent_handle_message")
		if err != nil {
			rc.Err(err)
			return
		}

		params := map[string]any{
			"id":   a.agent.Id,
			"name": msg.Name,
			"data": msg.Data,
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
				if str, ok := result.(string); ok {
					response.Data = &str
				} else {
					rc.Err(fmt.Errorf("unexpected result type: %T", result))
					return
				}
			}
			rc.Response(response)
		}

	default:
		rc.Unhandled()
	}
}

func (a *WasmAgentActor) startAgent(ctx context.Context) error {

	fnInfo, err := a.host.GetFunctionInfo("_modus_start_agent")
	if err != nil {
		return err
	}

	params := map[string]any{
		"id": a.agent.Id,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)

	_ = execInfo // TODO
	return err
}

func (a *WasmAgentActor) stopAgent(ctx context.Context) error {

	fnInfo, err := a.host.GetFunctionInfo("_modus_stop_agent")
	if err != nil {
		return err
	}

	params := map[string]any{
		"id": a.agent.Id,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)

	_ = execInfo // TODO
	return err
}

func (a *WasmAgentActor) getAgentState(ctx context.Context) (*string, error) {

	fnInfo, err := a.host.GetFunctionInfo("_modus_get_agent_state")
	if err != nil {
		return nil, err
	}

	params := map[string]any{
		"id": a.agent.Id,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
	if err != nil {
		return nil, err
	}

	result := execInfo.Result()
	if result == nil {
		return nil, nil
	}

	state, ok := result.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
	return &state, nil
}

func (a *WasmAgentActor) setAgentState(ctx context.Context, data string) error {
	fnInfo, err := a.host.GetFunctionInfo("_modus_set_agent_state")
	if err != nil {
		return err
	}

	params := map[string]any{
		"id":   a.agent.Id,
		"data": data,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
	if err != nil {
		return err
	}

	_ = execInfo // TODO
	return err
}
