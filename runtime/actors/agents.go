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

func StartAgent(ctx context.Context, agentName string) (*AgentInfo, error) {
	plugin, ok := plugins.GetPluginFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no plugin found in context")
	}

	agentId := xid.New().String()
	host := wasmhost.GetWasmHost(ctx)
	spawnActorForAgentAsync(host, plugin, agentId, agentName, false, nil)

	info := &AgentInfo{
		Id:     agentId,
		Name:   agentName,
		Status: AgentStatusStarting,
	}

	return info, nil
}

func spawnActorForAgentAsync(host wasmhost.WasmHost, plugin *plugins.Plugin, agentId, agentName string, resuming bool, initialState *string) {
	// We spawn the actor in a goroutine to avoid blocking while the actor is being spawned.
	// This allows many agents to be spawned in parallel, if needed.
	// Errors are logged but not returned, as the actor system will handle them.
	go func() {
		_, _ = spawnActorForAgent(host, plugin, agentId, agentName, resuming, initialState)
	}()
}

func spawnActorForAgent(host wasmhost.WasmHost, plugin *plugins.Plugin, agentId, agentName string, resuming bool, initialState *string) (*goakt.PID, error) {
	// The actor needs to spawn in its own context, so we don't pass one in to this function.
	// If we did, then when the original context was cancelled or completed, the actor initialization would be cancelled too.
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.WasmHostContextKey, host)
	ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, agentName)

	actor := newWasmAgentActor(agentId, agentName, host, plugin)
	actorName := getActorName(agentId)

	// note, don't call updateStatus here, because we can't publish an event yet
	if resuming {
		actor.status = AgentStatusResuming
		actor.initialState = initialState
	} else {
		actor.status = AgentStatusStarting
	}

	pid, err := _actorSystem.Spawn(ctx, actorName, actor)
	if err != nil {
		logger.Err(ctx, err).Msg("Error spawning actor for agent.")
	}
	return pid, err
}

func StopAgent(ctx context.Context, agentId string) (*AgentInfo, error) {
	pid, err := getActorPid(ctx, agentId)
	if err != nil {
		// see if it's in the database before erroring
		if agent, e := db.GetAgentState(ctx, agentId); e == nil {
			return &AgentInfo{
				Id:     agent.Id,
				Name:   agent.Name,
				Status: agent.Status,
			}, nil
		}

		return nil, fmt.Errorf("error stopping agent %s: %w", agentId, err)
	}

	// it was found, so we can stop it
	actor := pid.Actor().(*wasmAgentActor)
	actor.updateStatus(ctx, AgentStatusStopping)
	if err := pid.Shutdown(ctx); err != nil {
		return nil, fmt.Errorf("error stopping agent %s: %w", agentId, err)
	}

	return &AgentInfo{
		Id:     actor.agentId,
		Name:   actor.agentName,
		Status: actor.status,
	}, nil
}

func GetAgentInfo(ctx context.Context, agentId string) (*AgentInfo, error) {

	// Try the local actor system first.
	if pid, err := getActorPid(ctx, agentId); err == nil {
		actor := pid.Actor().(*wasmAgentActor)
		return &AgentInfo{
			Id:     actor.agentId,
			Name:   actor.agentName,
			Status: actor.status,
		}, nil
	}

	// Check the database as a fallback.
	// This is useful if the actor is terminated, or running on another node.
	if agent, err := db.GetAgentState(ctx, agentId); err == nil {
		return &AgentInfo{
			Id:     agent.Id,
			Name:   agent.Name,
			Status: agent.Status,
		}, nil
	}

	return nil, fmt.Errorf("agent %s not found", agentId)
}

type agentMessageResponse struct {
	Data  *string
	Error *string
}

func SendAgentMessage(ctx context.Context, agentId string, msgName string, data *string, timeout int64) (*agentMessageResponse, error) {

	pid, err := getActorPid(ctx, agentId)
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

	pid, err := getActorPid(ctx, agentId)
	if err != nil {
		return err
	}

	topicActor := _actorSystem.TopicActor()
	return pid.Tell(ctx, topicActor, pubMsg)
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

	if errors.Is(err, goakt.ErrActorNotFound) {
		// the actor might have been suspended or terminated, so check the database
		if agent, dbErr := db.GetAgentState(ctx, agentId); dbErr == nil {
			switch agent.Status {

			case AgentStatusSuspended:
				// the actor is suspended, so we can resume it (if we're in the right context)
				if host, ok := wasmhost.TryGetWasmHost(ctx); ok {
					plugin, ok := plugins.GetPluginFromContext(ctx)
					if !ok {
						return nil, fmt.Errorf("no plugin found in context for agent %s", agentId)
					}
					return spawnActorForAgent(host, plugin, agentId, agent.Name, true, &agent.Data)
				}

				return nil, fmt.Errorf("agent %s is suspended", agentId)

			case AgentStatusTerminated:
				return nil, fmt.Errorf("agent %s is terminated", agentId)

			default:
				// try one more time in case the actor was activated just after the initial check
				if _, pid, err := _actorSystem.ActorOf(ctx, actorName); err == nil {
					return pid, nil
				}

				// this means the agent is running on another node - TODO: handle this somehow
				return nil, fmt.Errorf("agent %s is %s, but not found in local actor system", agentId, agent.Status)
			}
		}
		return nil, fmt.Errorf("agent %s not found", agentId)
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
