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
	"errors"
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/messages"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	"github.com/rs/xid"
	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/goaktpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AgentInfo struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Status AgentStatus `json:"status"`
}

type AgentStatus = string

const (
	AgentStatusStarting   AgentStatus = "starting"
	AgentStatusRunning    AgentStatus = "running"
	AgentStatusSuspending AgentStatus = "suspending"
	AgentStatusSuspended  AgentStatus = "suspended"
	AgentStatusResuming   AgentStatus = "resuming"
	AgentStatusStopping   AgentStatus = "stopping"
	AgentStatusTerminated AgentStatus = "terminated"
)

const agentStatusEventName = "agentStatusUpdated"

type agentEventAction = string

const (
	agentEventActionInitialize agentEventAction = "initialize"
	agentEventActionSuspend    agentEventAction = "suspend"
	agentEventActionResume     agentEventAction = "resume"
	agentEventActionTerminate  agentEventAction = "terminate"
)

func StartAgent(ctx context.Context, agentName string) (*AgentInfo, error) {
	plugin, ok := plugins.GetPluginFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no plugin found in context")
	}

	agentId := xid.New().String()
	host := wasmhost.GetWasmHost(ctx)
	spawnActorForAgentAsync(host, plugin, agentId, agentName, true)

	info := &AgentInfo{
		Id:     agentId,
		Name:   agentName,
		Status: AgentStatusStarting,
	}

	return info, nil
}

func spawnActorForAgentAsync(host wasmhost.WasmHost, plugin *plugins.Plugin, agentId, agentName string, initializing bool) {
	// We spawn the actor in a goroutine to avoid blocking while the actor is being spawned.
	// This allows many agents to be spawned in parallel, if needed.
	// Errors are logged but not returned, as the actor system will handle them.
	go func() {
		_, _ = spawnActorForAgent(host, plugin, agentId, agentName, initializing)
	}()
}

func spawnActorForAgent(host wasmhost.WasmHost, plugin *plugins.Plugin, agentId, agentName string, initializing bool) (*goakt.PID, error) {
	// The actor needs to spawn in its own context, so we don't pass one in to this function.
	// If we did, then when the original context was cancelled or completed, the actor initialization would be cancelled too.
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.WasmHostContextKey, host)
	ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, agentName)

	actor := &wasmAgentActor{
		agentId:      agentId,
		agentName:    agentName,
		host:         host,
		plugin:       plugin,
		initializing: initializing,
	}

	actorName := getActorName(agentId)
	pid, err := _actorSystem.Spawn(ctx, actorName, actor)
	if err != nil {
		logger.Err(ctx, err).Msg("Error spawning actor for agent.")
	}
	return pid, err
}

func StopAgent(ctx context.Context, agentId string) (*AgentInfo, error) {
	info, pid, err := ensureAgentReady(ctx, agentId)
	if pid == nil && info != nil {
		return info, nil
	} else if err != nil {
		return nil, err
	}

	// shut down the actor, which will then stop the agent
	actor := pid.Actor().(*wasmAgentActor)
	if actor.status != AgentStatusStopping && actor.status != AgentStatusTerminated {
		actor.terminating = true
		if err := pid.Shutdown(ctx); err != nil {
			return nil, fmt.Errorf("error stopping agent %s: %w", agentId, err)
		}
	}

	return &AgentInfo{
		Id:     actor.agentId,
		Name:   actor.agentName,
		Status: actor.status,
	}, nil
}

func GetAgentInfo(ctx context.Context, agentId string) (*AgentInfo, error) {
	info, _, err := getAgentInfo(ctx, agentId)
	return info, err
}

func getAgentInfo(ctx context.Context, agentId string) (*AgentInfo, *goakt.PID, error) {
	pid, err := getActorPid(ctx, agentId)
	if errors.Is(err, goakt.ErrActorNotFound) {
		if agent, e := db.GetAgentState(ctx, agentId); e == nil {
			return &AgentInfo{
				Id:     agent.Id,
				Name:   agent.Name,
				Status: agent.Status,
			}, nil, nil
		}
		return nil, nil, fmt.Errorf("agent %s not found", agentId)
	} else if err != nil {
		return nil, nil, err
	}

	actor := pid.Actor().(*wasmAgentActor)
	return &AgentInfo{
		Id:     actor.agentId,
		Name:   actor.agentName,
		Status: actor.status,
	}, pid, nil
}

type agentMessageResponse struct {
	Data  *string
	Error *string
}

func ensureAgentReady(ctx context.Context, agentId string) (*AgentInfo, *goakt.PID, error) {
	const maxRetries = 8
	const baseDelay = 25 * time.Millisecond
	const maxDelay = 500 * time.Millisecond

	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		info, pid, err := getAgentInfo(ctx, agentId)

		if pid != nil {
			actor := pid.Actor().(*wasmAgentActor)
			if actor.status == AgentStatusRunning {
				return info, pid, nil
			}
			lastErr = fmt.Errorf("agent %s not ready (status: %s)", agentId, actor.status)
		} else if info != nil {
			switch info.Status {
			case AgentStatusStarting:
				lastErr = fmt.Errorf("agent %s is still starting", agentId)

			case AgentStatusSuspended:
				// Try to resume the agent (only on first attempt to avoid spam)
				if attempt == 0 {
					host := wasmhost.GetWasmHost(ctx)
					plugin, ok := plugins.GetPluginFromContext(ctx)
					if !ok {
						return info, nil, fmt.Errorf("no plugin found in context for agent %s", agentId)
					}

					// Resume asynchronously
					go func() {
						if _, err := spawnActorForAgent(host, plugin, agentId, info.Name, false); err != nil {
							logger.Err(context.Background(), err).Msgf("Failed to resume agent %s", agentId)
						}
					}()
				}
				lastErr = fmt.Errorf("agent %s is resuming", agentId)

			case AgentStatusTerminated:
				// Permanent failure - don't retry
				return info, nil, fmt.Errorf("agent %s is terminated", agentId)

			default:
				// Other statuses - don't retry
				return info, nil, fmt.Errorf("agent %s is %s, but not found in local actor system", agentId, info.Status)
			}
		} else if err != nil {
			// Handle the error from getAgentInfo
			lastErr = err
		} else {
			// Agent doesn't exist at all - permanent failure
			return nil, nil, fmt.Errorf("agent %s not found", agentId)
		}

		// Exit early if this is the last attempt
		if attempt == maxRetries-1 {
			break
		}

		// Wait before retrying with exponential backoff
		delay := time.Duration(float64(baseDelay) * (1.5*float64(attempt) + 1))
		if delay > maxDelay {
			delay = maxDelay
		}

		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	// All retries exhausted
	return nil, nil, fmt.Errorf("agent %s not ready after %d attempts: %v", agentId, maxRetries, lastErr)
}

func SendAgentMessage(ctx context.Context, agentId string, msgName string, data *string, timeout int64) (*agentMessageResponse, error) {

	_, pid, err := ensureAgentReady(ctx, agentId)
	if err != nil {
		e := err.Error()
		return &agentMessageResponse{Error: &e}, err
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
		return &agentMessageResponse{Data: response.Data}, nil
	} else {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

func PublishAgentEvent(ctx context.Context, agentId, eventName string, eventData *string) error {

	var data any
	if eventData != nil {
		if err := utils.JsonDeserialize([]byte(*eventData), &data); err != nil {
			return fmt.Errorf("error deserializing event data: %w", err)
		}
	}

	dataValue, err := structpb.NewValue(data)
	if err != nil {
		return fmt.Errorf("error creating event data value: %w", err)
	}

	event := &messages.AgentEventMessage{
		Name:      eventName,
		Data:      dataValue,
		Timestamp: timestamppb.Now(),
	}

	eventMsg, err := anypb.New(event)
	if err != nil {
		return fmt.Errorf("error creating event message: %w", err)
	}

	pubMsg := &goaktpb.Publish{
		Id:      xid.New().String(),
		Topic:   getAgentTopic(agentId),
		Message: eventMsg,
	}

	topicActor := _actorSystem.TopicActor()

	if pid, err := getActorPid(ctx, agentId); err == nil {
		return pid.Tell(ctx, topicActor, pubMsg)
	}

	// publish anonymously if the actor is not found
	if errors.Is(err, goakt.ErrActorNotFound) {
		// For now, we use the topic actor directly to publish the message.
		// See https://github.com/Tochemey/goakt/pull/761
		// TODO: use goakt.Tell after it's fixed
		return topicActor.Tell(ctx, topicActor, pubMsg)
	}

	return err
}

func getActorName(agentId string) string {
	return "agent-" + agentId
}

func getAgentTopic(agentId string) string {
	return getActorName(agentId) + ".events"
}

func getActorPid(ctx context.Context, agentId string) (*goakt.PID, error) {
	if _, err := xid.FromString(agentId); err != nil {
		return nil, fmt.Errorf("invalid agent ID format: %s", agentId)
	}

	actorName := getActorName(agentId)
	_, pid, err := _actorSystem.ActorOf(ctx, actorName)
	if err == nil {
		return pid, nil
	}

	return nil, fmt.Errorf("error getting actor for agent %s: %w", agentId, err)
}

func ListActiveAgents(ctx context.Context) ([]AgentInfo, error) {
	agents, err := db.QueryActiveAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing active agents: %w", err)
	}

	results := make([]AgentInfo, 0, len(agents))
	for _, agent := range agents {
		results = append(results, AgentInfo{
			Id:     agent.Id,
			Name:   agent.Name,
			Status: agent.Status,
		})
	}

	return results, nil
}

func ListLocalAgents() []AgentInfo {
	if _actorSystem == nil {
		return nil
	}

	actors := _actorSystem.Actors()
	results := make([]AgentInfo, 0, len(actors))

	for _, pid := range actors {
		if actor, ok := pid.Actor().(*wasmAgentActor); ok {
			results = append(results, AgentInfo{
				Id:     actor.agentId,
				Name:   actor.agentName,
				Status: actor.status,
			})
		}
	}

	return results
}
