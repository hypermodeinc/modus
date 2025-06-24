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
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/rs/xid"
	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/goaktpb"
)

type agentEvent struct {
	Name      string `json:"name"`
	Data      any    `json:"data"`
	Timestamp string `json:"timestamp"`
}

func SubscribeForAgentEvents(ctx context.Context, agentId string, update func(data []byte), done func()) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	// Go directly to the database for the agent status, because we don't want subscribing to event to fail
	// if there is any issue with the agent actor.
	if a, err := getAgentInfoFromDatabase(ctx, agentId); err != nil {
		return err
	} else if a.Status == AgentStatusStopping || a.Status == AgentStatusTerminated {
		return fmt.Errorf("agent %s is %s, cannot subscribe to events", agentId, a.Status)
	}

	if update == nil {
		update = func(data []byte) {}
	}
	if done == nil {
		done = func() {}
	}
	actor := &subscriptionActor{
		update: update,
		done:   done,
	}

	// Spawn a subscription actor that is bound to the graphql subscription on this node.
	// It needs to be long-lived, because it will need to stay alive as long as the client is connected.
	// It cannot be relocated to another node, because it is bound to http request of the GraphQL subscription on this node.
	// It cannot be spawned as a child of the agent actor, because the subscription needs to be maintained even if the agent actor is suspended.

	actorName := "subscription-" + xid.New().String()
	subActor, err := _actorSystem.Spawn(ctx, actorName, actor, goakt.WithLongLived(), goakt.WithRelocationDisabled())
	if err != nil {
		return fmt.Errorf("failed to spawn subscription actor: %w", err)
	}

	topic := getAgentTopic(agentId)
	subscribe := &goaktpb.Subscribe{Topic: topic}

	if err := subActor.Tell(ctx, _actorSystem.TopicActor(), subscribe); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	// When the context is done, we will unsubscribe and stop the subscription actor.
	// For example, the GraphQL subscription is closed or the client disconnects.
	go func() {
		<-ctx.Done()

		// reset cancellation because the original context is already done
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()

		logger.Debug(ctx).Msgf("Unsubscribing from topic %s and shutting down subscription actor %s", topic, subActor.Name())

		unsubscribe := &goaktpb.Unsubscribe{Topic: topic}
		if err := subActor.Tell(ctx, _actorSystem.TopicActor(), unsubscribe); err != nil {
			logger.Err(ctx, err).Msg("Failed to unsubscribe from topic")
		}

		if err := subActor.Shutdown(ctx); err != nil {
			logger.Err(ctx, err).Msg("Failed to shut down subscription actor")
		}
	}()

	return nil
}

type subscriptionActor struct {
	update func(data []byte)
	done   func()
}

func (a *subscriptionActor) PreStart(ac *goakt.Context) error {
	return nil
}

func (a *subscriptionActor) PostStop(ac *goakt.Context) error {
	a.done()
	return nil
}

func (a *subscriptionActor) Receive(rc *goakt.ReceiveContext) {
	if msg, ok := rc.Message().(*messages.AgentEvent); ok {
		span, _ := utils.NewSentrySpanForCurrentFunc(rc.Context())
		defer span.Finish()

		event := &agentEvent{
			Name:      msg.Name,
			Data:      msg.Data,
			Timestamp: msg.Timestamp.AsTime().Format(utils.TimeFormat),
		}
		if data, err := utils.JsonSerialize(event); err != nil {
			rc.Err(fmt.Errorf("failed to serialize agent event message: %w", err))
		} else {
			a.update(data)
		}

		if msg.Name == agentStatusEventName {
			status := AgentStatus(msg.Data.GetStructValue().Fields["status"].GetStringValue())
			if status == AgentStatusTerminated {
				rc.Shutdown()
			}
		}

		return
	}

	rc.Unhandled()
}
