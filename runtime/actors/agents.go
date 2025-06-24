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

	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/goaktpb"

	"github.com/rs/xid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type pidContextKey struct{}

type AgentInfo struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Status AgentStatus `json:"status"`
}

type AgentStatus string

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

type agentEventAction string

const (
	agentEventActionInitialize agentEventAction = "initialize"
	agentEventActionSuspend    agentEventAction = "suspend"
	agentEventActionResume     agentEventAction = "resume"
	agentEventActionTerminate  agentEventAction = "terminate"
)

func StartAgent(ctx context.Context, agentName string) (*AgentInfo, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	plugin, ok := plugins.GetPluginFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no plugin found in context")
	}

	agentId := xid.New().String()
	if err := spawnActorForAgent(ctx, plugin.Name(), agentId, agentName, true); err != nil {
		return nil, fmt.Errorf("error spawning actor for agent %s: %w", agentId, err)
	}

	return GetAgentInfo(ctx, agentId)
}

func spawnActorForAgent(ctx context.Context, pluginName, agentId, agentName string, initializing bool) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	ctx = context.WithoutCancel(ctx)
	ctx = context.WithValue(ctx, utils.AgentIdContextKey, agentId)
	ctx = context.WithValue(ctx, utils.AgentNameContextKey, agentName)

	actor := &wasmAgentActor{
		// this only works because we always spawn locally the first time
		initializing: initializing,
	}

	actorName := getActorName(agentId)
	_, err := _actorSystem.Spawn(ctx, actorName, actor,
		goakt.WithLongLived(),
		goakt.WithDependencies(&wasmAgentInfo{
			AgentName:  agentName,
			PluginName: pluginName,
		}),
	)

	// Important: Wait for the actor system to sync with the cluster before proceeding.
	// This ensures consistency across the cluster, so we don't accidentally spawn the same actor multiple times.
	// GoAkt does not resolve such inconsistencies automatically, so we need to handle this manually.
	// A short sync time should not be noticeable by the user.
	waitForClusterSync(ctx)

	return err
}

func StopAgent(ctx context.Context, agentId string) (*AgentInfo, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	actorName := getActorName(agentId)
	if err := tell(ctx, actorName, &messages.ShutdownAgent{}); err != nil {
		if !errors.Is(err, goakt.ErrActorNotFound) {
			return nil, fmt.Errorf("error stopping agent %s: %w", agentId, err)
		}
	}

	// Don't ask the actor, because it might already be stopped.
	info, err := getAgentInfoFromDatabase(ctx, agentId)
	if err != nil {
		return nil, err
	}

	// If the agent is not yet terminated, we'll send back "stopping"
	// so we don't have to wait for the actor to be stopped synchronously.
	if info.Status != AgentStatusTerminated {
		info.Status = AgentStatusStopping
	}
	return info, nil

}

func getAgentInfoFromDatabase(ctx context.Context, agentId string) (*AgentInfo, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	if agent, e := db.GetAgentState(ctx, agentId); e == nil {
		return &AgentInfo{
			Id:     agent.Id,
			Name:   agent.Name,
			Status: AgentStatus(agent.Status),
		}, nil
	}
	return nil, fmt.Errorf("agent %s not found", agentId)
}

func GetAgentInfo(ctx context.Context, agentId string) (*AgentInfo, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	actorName := getActorName(agentId)
	request := &messages.AgentInfoRequest{}

	// We first try to ask the actor for its info.  Use a short timeout to avoid blocking indefinitely.
	response, err := ask(ctx, actorName, request, 500*time.Millisecond)
	if err == nil {
		msg := response.(*messages.AgentInfoResponse)
		return &AgentInfo{
			Id:     agentId,
			Name:   msg.Name,
			Status: AgentStatus(msg.Status),
		}, nil
	}

	// If the actor is not responding, we can check the database for the agent state.
	// This is useful for agents that are still starting, are terminated or suspended, or just busy processing another request.
	if errors.Is(err, goakt.ErrActorNotFound) || errors.Is(err, goakt.ErrRequestTimeout) || errors.Is(err, goakt.ErrDead) || errors.Is(err, goakt.ErrRemoteSendFailure) {
		return getAgentInfoFromDatabase(ctx, agentId)
	}

	return nil, fmt.Errorf("error getting agent info: %w", err)
}

type agentMessageResponse struct {
	Data  *string
	Error *string
}

func newAgentMessageDataResponse(data *string) *agentMessageResponse {
	return &agentMessageResponse{Data: data}
}

func newAgentMessageErrorResponse(errMsg string) *agentMessageResponse {
	return &agentMessageResponse{Error: &errMsg}
}

func SendAgentMessage(ctx context.Context, agentId string, msgName string, data *string, timeout int64) (*agentMessageResponse, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	actorName := getActorName(agentId)

	msg := &messages.AgentRequest{
		Name:    msgName,
		Data:    data,
		Respond: timeout > 0,
	}

	var err error
	var res proto.Message
	if timeout == 0 {
		err = tell(ctx, actorName, msg)
	} else {
		res, err = ask(ctx, actorName, msg, time.Duration(timeout))
	}

	if errors.Is(err, goakt.ErrActorNotFound) {
		return newAgentMessageErrorResponse("agent not found"), nil
	} else if err != nil {
		return nil, fmt.Errorf("error sending message to agent: %w", err)
	}

	if res == nil {
		return newAgentMessageDataResponse(nil), nil
	} else if response, ok := res.(*messages.AgentResponse); ok {
		return newAgentMessageDataResponse(response.Data), nil
	} else {
		return nil, fmt.Errorf("unexpected agent response type: %T", res)
	}
}

func PublishAgentEvent(ctx context.Context, agentId, eventName string, eventData *string) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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

	event := &messages.AgentEvent{
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

	// if the pid is in context, we're being called as a host function
	if pid, ok := ctx.Value(pidContextKey{}).(*goakt.PID); ok {
		return pid.Tell(ctx, topicActor, pubMsg)
	}

	// otherwise, we try to get the actor PID for the agent (we should avoid this)
	pid, err := _actorSystem.LocalActor(getActorName(agentId))
	if err == nil {
		logger.Warn(ctx).Str("event", eventName).Any("data", eventData).Msg("Agent actor not in context. Using lookup to publish event.")
		return pid.Tell(ctx, topicActor, pubMsg)
	}

	// publish anonymously if the actor is not found (we should avoid this)
	if errors.Is(err, goakt.ErrActorNotFound) {
		logger.Warn(ctx).Str("event", eventName).Any("data", eventData).Msg("Agent actor not found. Publishing event anonymously.")
		return goakt.Tell(ctx, topicActor, pubMsg)
	}

	return err
}

func getActorName(agentId string) string {
	return "agent-" + agentId
}

func getAgentTopic(agentId string) string {
	return getActorName(agentId) + ".events"
}

func ListActiveAgents(ctx context.Context) ([]AgentInfo, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	agents, err := db.QueryActiveAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing active agents: %w", err)
	}

	results := make([]AgentInfo, 0, len(agents))
	for _, agent := range agents {
		results = append(results, AgentInfo{
			Id:     agent.Id,
			Name:   agent.Name,
			Status: AgentStatus(agent.Status),
		})
	}

	return results, nil
}

func ListLocalAgents(ctx context.Context) []AgentInfo {
	if _actorSystem == nil {
		return nil
	}

	span, _ := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
