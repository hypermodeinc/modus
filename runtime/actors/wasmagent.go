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

	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/messages"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	wasm "github.com/tetratelabs/wazero/api"
	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/goaktpb"
)

type wasmAgentActor struct {
	agentId      string
	agentName    string
	status       AgentStatus
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
	case AgentStatusStarting:
		logger.Info(ctx).Msg("Starting agent.")
	case AgentStatusResuming:
		logger.Info(ctx).Msg("Resuming agent.")
	case AgentStatusSuspended:
		// note, don't call updateStatus here, because we can't publish an event yet
		a.status = AgentStatusResuming
		logger.Info(ctx).Msg("Resuming agent.")
	default:
		return fmt.Errorf("invalid agent status for actor PreStart: %s", a.status)
	}

	if err := a.saveState(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error saving agent state.")
	}

	a.buffers = utils.NewOutputBuffers()
	if mod, err := a.host.GetModuleInstance(ctx, a.plugin, a.buffers); err != nil {
		return err
	} else {
		a.module = mod
	}

	return nil
}

func (a *wasmAgentActor) PostStop(ac *goakt.Context) error {
	ctx := a.newContext()
	defer a.module.Close(ctx)

	switch a.status {
	case AgentStatusRunning:
		a.updateStatus(ctx, AgentStatusSuspending)
		logger.Info(ctx).Msg("Suspending agent.")
	case AgentStatusSuspending:
		logger.Info(ctx).Msg("Suspending agent.")
	case AgentStatusStopping:
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
	case AgentStatusSuspending:
		a.updateStatus(ctx, AgentStatusSuspended)
		if err := a.saveState(ctx); err != nil {
			return err
		}
		logger.Info(ctx).Msg("Agent suspended successfully.")
	case AgentStatusStopping:
		a.updateStatus(ctx, AgentStatusTerminated)
		if err := a.saveState(ctx); err != nil {
			return err
		}
		logger.Info(ctx).Dur("duration_ms", duration).Msg("Agent terminated successfully.")
	default:
		return fmt.Errorf("invalid agent status for actor PostStop: %s", a.status)
	}

	return nil
}

func (a *wasmAgentActor) updateStatus(ctx context.Context, status AgentStatus) {
	// set the agent status (this is the important part)
	a.status = status

	// try to publish an event, but it's ok if it fails
	data := fmt.Sprintf(`{"status":"%s"}`, a.status)
	_ = PublishAgentEvent(ctx, a.agentId, agentStatusEventName, &data)
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
	case *goaktpb.PostStart:

		if err := a.activateAgent(ctx); err != nil {
			logger.Err(ctx, err).Msg("Error activating agent.")
			rc.Err(err)
			return
		}

		if a.status == AgentStatusResuming {
			if err := a.setAgentState(ctx, a.initialState); err != nil {
				logger.Warn(ctx).Err(err).Msg("Unable to restore agent state.")
			}
			a.initialState = nil
		}

		if a.status == AgentStatusResuming {
			logger.Info(ctx).Msg("Agent resumed successfully.")
		} else {
			logger.Info(ctx).Msg("Agent started successfully.")
		}

		a.updateStatus(ctx, AgentStatusRunning)
		if err := a.saveState(ctx); err != nil {
			logger.Err(ctx, err).Msg("Error saving agent state.")
			rc.Err(err)
			return
		}

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
		"reloading": a.status == AgentStatusResuming,
	}

	_, err = a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
	return err
}

func (a *wasmAgentActor) shutdownAgent(ctx context.Context) error {

	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_shutdown")
	if err != nil {
		return err
	}

	params := map[string]any{
		"suspending": a.status == AgentStatusSuspending,
	}

	_, err = a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
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

	_, err = a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
	return err
}

func (a *wasmAgentActor) reloadModule(ctx context.Context, plugin *plugins.Plugin) error {

	// the context may not have these values set
	ctx = context.WithValue(ctx, utils.PluginContextKey, a.plugin)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, a.agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, a.agentName)

	logger.Info(ctx).Msg("Reloading module for agent.")

	a.updateStatus(ctx, AgentStatusSuspending)
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
	a.updateStatus(ctx, AgentStatusSuspended)

	// create a new module instance and assign it to the actor
	a.plugin = plugin
	a.buffers = utils.NewOutputBuffers()
	mod, err := a.host.GetModuleInstance(ctx, a.plugin, a.buffers)
	if err != nil {
		return err
	}
	a.module = mod

	// activate the agent in the new module instance
	a.updateStatus(ctx, AgentStatusResuming)
	if err := a.activateAgent(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error reloading agent.")
		return err
	}

	// restore the state in the new module instance
	if err := a.setAgentState(ctx, state); err != nil {
		logger.Err(ctx, err).Msg("Error setting agent state.")
		return err
	}

	a.updateStatus(ctx, AgentStatusRunning)
	logger.Info(ctx).Msg("Agent reloaded module successfully.")

	return nil
}
