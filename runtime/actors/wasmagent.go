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
)

type wasmAgentActor struct {
	agentId      string
	agentName    string
	status       AgentStatus
	plugin       *plugins.Plugin
	host         wasmhost.WasmHost
	module       wasm.Module
	buffers      utils.OutputBuffers
	initializing bool
	terminating  bool
}

func (a *wasmAgentActor) PreStart(ac *goakt.Context) error {
	ctx := a.newContext()

	if err := a.activateAgent(ctx); err != nil {
		return fmt.Errorf("error activating agent: %w", err)
	}

	if a.initializing {
		if err := a.startAgent(ctx); err != nil {
			return fmt.Errorf("error starting agent: %w", err)
		}
		a.initializing = false
	} else {
		if err := a.resumeAgent(ctx); err != nil {
			return fmt.Errorf("error resuming agent: %w", err)
		}
	}

	return nil
}

func (a *wasmAgentActor) PostStop(ac *goakt.Context) error {
	ctx := a.newContext()
	defer a.module.Close(ctx)

	if a.terminating {
		if err := a.stopAgent(ctx); err != nil {
			return fmt.Errorf("error stopping agent: %w", err)
		}
	} else {
		if err := a.suspendAgent(ctx); err != nil {
			return fmt.Errorf("error suspending agent: %w", err)
		}
	}

	return nil
}

func (a *wasmAgentActor) updateStatus(ctx context.Context, status AgentStatus) error {
	if a.status == status {
		return nil
	}

	a.status = status

	if err := db.UpdateAgentStatus(ctx, a.agentId, status); err != nil {
		return fmt.Errorf("error updating agent status in database: %w", err)
	}

	data := fmt.Sprintf(`{"status":"%s"}`, a.status)
	if err := PublishAgentEvent(ctx, a.agentId, agentStatusEventName, &data); err != nil {
		return fmt.Errorf("error publishing agent status event: %w", err)
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
		UpdatedAt: time.Now().UTC().Format(utils.TimeFormat),
	}); err != nil {
		return fmt.Errorf("error saving agent state to database: %w", err)
	}

	return nil
}

func (a *wasmAgentActor) restoreState(ctx context.Context) error {
	if a.module == nil {
		return fmt.Errorf("module is not initialized")
	}

	data, err := db.GetAgentState(ctx, a.agentId)
	if err != nil {
		return fmt.Errorf("error getting agent state from database: %w", err)
	}

	if err := a.setAgentState(ctx, &data.Data); err != nil {
		return fmt.Errorf("error setting agent state: %w", err)
	}

	return nil
}

func (a *wasmAgentActor) Receive(rc *goakt.ReceiveContext) {
	ctx := a.newContext()

	// NOTE: GoAkt will send a goaktpb.PostStart message, but we don't need to do anything with it.

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

	a.buffers = utils.NewOutputBuffers()
	if mod, err := a.host.GetModuleInstance(ctx, a.plugin, a.buffers); err != nil {
		return fmt.Errorf("error getting module instance in actor pre-start: %w", err)
	} else {
		a.module = mod
	}

	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_activate")
	if err != nil {
		return err
	}

	params := map[string]any{
		"name": a.agentName,
		"id":   a.agentId,
	}

	_, err = a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
	return err
}

func (a *wasmAgentActor) startAgent(ctx context.Context) error {
	logger.Info(ctx).Msg("Starting agent.")
	if err := a.saveState(ctx); err != nil {
		return err
	}
	if err := a.updateStatus(ctx, AgentStatusStarting); err != nil {
		return err
	}
	if err := a.callEventHandler(ctx, agentEventActionInitialize); err != nil {
		return err
	}
	if err := a.updateStatus(ctx, AgentStatusRunning); err != nil {
		return err
	}
	logger.Info(ctx).Msg("Agent started successfully.")
	return nil
}

func (a *wasmAgentActor) resumeAgent(ctx context.Context) error {
	logger.Info(ctx).Msg("Resuming agent.")
	if err := a.updateStatus(ctx, AgentStatusResuming); err != nil {
		return err
	}
	if err := a.restoreState(ctx); err != nil {
		return err
	}
	if err := a.callEventHandler(ctx, agentEventActionResume); err != nil {
		return err
	}
	if err := a.updateStatus(ctx, AgentStatusRunning); err != nil {
		return err
	}
	logger.Info(ctx).Msg("Agent resumed successfully.")
	return nil
}

func (a *wasmAgentActor) suspendAgent(ctx context.Context) error {
	logger.Info(ctx).Msg("Suspending agent.")
	if err := a.updateStatus(ctx, AgentStatusSuspending); err != nil {
		return err
	}
	if err := a.callEventHandler(ctx, agentEventActionSuspend); err != nil {
		return err
	}
	if err := a.saveState(ctx); err != nil {
		return err
	}
	if err := a.updateStatus(ctx, AgentStatusSuspended); err != nil {
		return err
	}
	logger.Info(ctx).Msg("Agent suspended successfully.")
	return nil
}

func (a *wasmAgentActor) stopAgent(ctx context.Context) error {
	logger.Info(ctx).Msg("Stopping agent.")
	if err := a.updateStatus(ctx, AgentStatusStopping); err != nil {
		return err
	}
	if err := a.callEventHandler(ctx, agentEventActionTerminate); err != nil {
		return err
	}
	if err := a.saveState(ctx); err != nil {
		return err
	}
	if err := a.updateStatus(ctx, AgentStatusTerminated); err != nil {
		return err
	}
	logger.Info(ctx).Msg("Agent stopped successfully.")
	return nil
}

func (a *wasmAgentActor) callEventHandler(ctx context.Context, action agentEventAction) error {
	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_handle_event")
	if err != nil {
		return err
	}

	params := map[string]any{
		"action": action,
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
