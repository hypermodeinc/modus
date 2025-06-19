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
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
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
	pluginName   string
	plugin       *plugins.Plugin
	host         wasmhost.WasmHost
	module       wasm.Module
	buffers      utils.OutputBuffers
	initializing bool
}

func (a *wasmAgentActor) PreStart(ac *goakt.Context) error {
	ctx := ac.Context()

	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	wasmExt := ac.Extension(wasmExtensionId).(*wasmExtension)
	a.host = wasmExt.host
	a.buffers = utils.NewOutputBuffers()

	if id, ok := strings.CutPrefix(ac.ActorName(), "agent-"); ok {
		a.agentId = id
	} else {
		return fmt.Errorf("actor name %s does not start with 'agent-'", ac.ActorName())
	}

	if info, err := getAgentInfo(ac); err != nil {
		return err
	} else {
		a.agentName = info.AgentName
		a.pluginName = info.PluginName
	}

	if err := a.activateAgent(ctx); err != nil {
		return fmt.Errorf("error activating agent: %w", err)
	}
	return nil
}

func (a *wasmAgentActor) Receive(rc *goakt.ReceiveContext) {
	ctx := a.augmentContext(rc.Context(), rc.Self())
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	switch msg := rc.Message().(type) {

	case *messages.AgentRequest:
		if err := a.handleAgentRequest(ctx, rc, msg); err != nil {
			rc.Err(err)
		}

	case *goaktpb.PostStart:
		if a.initializing {
			if err := a.startAgent(ctx); err != nil {
				rc.Err(fmt.Errorf("error starting agent: %w", err))
			}
			a.initializing = false
		} else {
			if err := a.resumeAgent(ctx); err != nil {
				rc.Err(fmt.Errorf("error resuming agent: %w", err))
			}
		}

	case *messages.AgentInfoRequest:
		rc.Response(&messages.AgentInfoResponse{
			Name:   a.agentName,
			Status: string(a.status),
		})

	case *messages.RestartAgent:
		if a.status != AgentStatusRunning {
			logger.Warn(ctx).Msgf("Agent is not %s, cannot restart now.", a.status)
			return
		}
		if err := a.suspendAgent(ctx); err != nil {
			rc.Err(fmt.Errorf("error suspending agent: %w", err))
			return
		}
		if err := a.deactivateAgent(ctx); err != nil {
			rc.Err(fmt.Errorf("error deactivating agent: %w", err))
			return
		}
		if err := a.activateAgent(ctx); err != nil {
			rc.Err(fmt.Errorf("error reactivating agent: %w", err))
			return
		}
		if err := a.resumeAgent(ctx); err != nil {
			rc.Err(fmt.Errorf("error resuming agent: %w", err))
			return
		}

	case *messages.ShutdownAgent:
		if a.status == AgentStatusStopping {
			logger.Warn(ctx).Msg("Agent is already stopping, cannot shutdown again.")
			return
		}
		if err := a.stopAgent(ctx); err != nil {
			rc.Err(fmt.Errorf("error stopping agent: %w", err))
		}
		rc.Shutdown()

	default:
		logger.Warn(ctx).Msgf("Unhandled message type %T in wasm agent actor.", msg)
		rc.Unhandled()
	}
}

func (a *wasmAgentActor) PostStop(ac *goakt.Context) error {
	ctx := ac.Context()
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	// suspend the agent if it's not already suspended or terminated
	if a.status != AgentStatusSuspended && a.status != AgentStatusTerminated {
		if err := a.suspendAgent(ctx); err != nil {
			logger.Err(ctx, err).Msg("Error suspending agent.")
			// don't return on error - we'll still try to deactivate the agent
		}
	}

	// deactivate the agent to clean up resources
	if err := a.deactivateAgent(ctx); err != nil {
		return fmt.Errorf("error deactivating agent: %w", err)
	}
	return nil
}

func (a *wasmAgentActor) handleAgentRequest(ctx context.Context, rc *goakt.ReceiveContext, msg *messages.AgentRequest) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	if a.status != AgentStatusRunning {
		return fmt.Errorf("cannot process message because agent is %s", a.status)
	}

	logger.Info(ctx).Str("msg_name", msg.Name).Msg("Received message.")
	start := time.Now()

	fnInfo, err := a.host.GetFunctionInfo("_modus_agent_handle_message")
	if err != nil {
		return err
	}

	params := map[string]any{
		"msgName": msg.Name,
		"data":    msg.Data,
	}

	execInfo, err := a.host.CallFunctionInModule(ctx, a.module, a.buffers, fnInfo, params)
	if err != nil {
		return err
	}

	if msg.Respond {
		result := execInfo.Result()
		response := &messages.AgentResponse{}
		if result != nil {
			switch result := result.(type) {
			case string:
				response.Data = &result
			case *string:
				response.Data = result
			default:
				err := fmt.Errorf("unexpected result type: %T", result)
				logger.Err(ctx, err).Msg("Error handling message.")
				return err
			}
		}
		// responding will cancel the context when finished, so we need to ensure
		// that the context used for the remainder of this function is not cancelled.
		ctx = context.WithoutCancel(ctx)
		rc.Response(response)
	}

	duration := time.Since(start)
	logger.Info(ctx).Str("msg_name", msg.Name).Dur("duration_ms", duration).Msg("Message handled successfully.")

	// save the state after handling the message to ensure the state is up to date in case of hard termination
	if err := a.saveState(ctx); err != nil {
		logger.Err(ctx, err).Msg("Error saving agent state.")
	}

	return nil
}

func (a *wasmAgentActor) updateStatus(ctx context.Context, status AgentStatus) error {
	if a.status == status {
		return nil
	}

	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	a.status = status

	if err := db.UpdateAgentStatus(ctx, a.agentId, string(status)); err != nil {
		return fmt.Errorf("error updating agent status in database: %w", err)
	}

	data := fmt.Sprintf(`{"status":"%s"}`, a.status)
	if err := PublishAgentEvent(ctx, a.agentId, agentStatusEventName, &data); err != nil {
		return fmt.Errorf("error publishing agent status event: %w", err)
	}

	return nil
}

func (a *wasmAgentActor) saveState(ctx context.Context) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
		Status:    string(a.status),
		Data:      data,
		UpdatedAt: time.Now(),
	}); err != nil {
		return fmt.Errorf("error saving agent state to database: %w", err)
	}

	return nil
}

func (a *wasmAgentActor) restoreState(ctx context.Context) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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

func (a *wasmAgentActor) augmentContext(ctx context.Context, pid *goakt.PID) context.Context {
	if ctx.Value(utils.WasmHostContextKey) == nil {
		ctx = context.WithValue(ctx, utils.WasmHostContextKey, a.host)
	}
	if ctx.Value(utils.PluginContextKey) == nil {
		ctx = context.WithValue(ctx, utils.PluginContextKey, a.plugin)
	}
	if ctx.Value(utils.AgentIdContextKey) == nil {
		ctx = context.WithValue(ctx, utils.AgentIdContextKey, a.agentId)
	}
	if ctx.Value(utils.AgentNameContextKey) == nil {
		ctx = context.WithValue(ctx, utils.AgentNameContextKey, a.agentName)
	}
	if ctx.Value(pidContextKey{}) == nil {
		ctx = context.WithValue(ctx, pidContextKey{}, pid)
	}
	return ctx
}

func (a *wasmAgentActor) activateAgent(ctx context.Context) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	if plugin, found := pluginmanager.GetPluginByName(a.pluginName); found {
		a.plugin = plugin
	} else {
		return fmt.Errorf("plugin %s not found", a.pluginName)
	}

	if mod, err := a.host.GetModuleInstance(ctx, a.plugin, a.buffers); err != nil {
		return err
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

func (a *wasmAgentActor) deactivateAgent(ctx context.Context) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	if err := a.module.Close(ctx); err != nil {
		return err
	}
	a.module = nil
	return nil
}

func (a *wasmAgentActor) startAgent(ctx context.Context) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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

const agentInfoDependencyId = "agent-info"

type wasmAgentInfo struct {
	AgentName  string `json:"agent"`
	PluginName string `json:"plugin"`
}

func (*wasmAgentInfo) ID() string {
	return agentInfoDependencyId
}

func (info *wasmAgentInfo) MarshalBinary() ([]byte, error) {
	jsonBytes, err := utils.JsonSerialize(info)
	if err != nil {
		return nil, fmt.Errorf("error serializing agent info to JSON: %w", err)
	}
	return jsonBytes, nil
}

func (info *wasmAgentInfo) UnmarshalBinary(data []byte) error {
	if err := utils.JsonDeserialize(data, info); err != nil {
		return fmt.Errorf("error deserializing agent info from JSON: %w", err)
	}
	return nil
}

func getAgentInfo(ac *goakt.Context) (*wasmAgentInfo, error) {
	deps := ac.Dependencies()
	for _, dep := range deps {
		if dep.ID() == agentInfoDependencyId {
			return dep.(*wasmAgentInfo), nil
		}
	}
	return nil, fmt.Errorf("agent info dependency not found")
}
