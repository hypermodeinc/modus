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
	"time"

	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/messages"
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	goakt "github.com/tochemey/goakt/v3/actor"
)

const actorSystemName = "modus"
const defaultAskTimeout = 10 * time.Second

var _actorSystem goakt.ActorSystem
var _remoting *goakt.Remoting

func Initialize(ctx context.Context) {

	opts := []goakt.Option{
		goakt.WithLogger(newActorLogger(logger.Get(ctx))),
		goakt.WithPubSub(),
		goakt.WithActorInitTimeout(10 * time.Second), // TODO: adjust this value, or make it configurable
		goakt.WithActorInitMaxRetries(1),             // TODO: adjust this value, or make it configurable

		// For now, keep passivation disabled so that agents can perform long-running tasks without the actor stopping.
		// TODO: Revisit this after https://github.com/Tochemey/goakt/issues/764 is resolved.
		goakt.WithPassivationDisabled(),
	}
	opts = append(opts, clusterOptions()...)

	if actorSystem, err := goakt.NewActorSystem(actorSystemName, opts...); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to create actor system.")
	} else if err := actorSystem.Start(ctx); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to start actor system.")
	} else {
		_actorSystem = actorSystem
	}

	_remoting = goakt.NewRemoting()

	logger.Info(ctx).Msg("Actor system started.")

	pluginmanager.RegisterPluginLoadedCallback(loadAgentActors)
}

func loadAgentActors(ctx context.Context, plugin *plugins.Plugin) error {
	// restart local actors that are already running, giving them the new plugin instance
	actors := _actorSystem.Actors()
	runningAgents := make(map[string]bool, len(actors))
	for _, pid := range actors {
		if a, ok := pid.Actor().(*wasmAgentActor); ok {
			runningAgents[a.agentId] = true
			a.plugin = plugin
			if err := goakt.Tell(ctx, pid, &messages.RestartAgent{}); err != nil {
				logger.Err(ctx, err).Str("agent_id", a.agentId).Msg("Failed to send restart agent message to actor.")
			}
		}
	}

	// spawn actors for agents with state in the database, that are not already running
	// TODO: when we scale out to allow more nodes in the cluster, we'll need to decide
	// which node is responsible for spawning each actor.
	agents, err := db.QueryActiveAgents(ctx)
	if err != nil {
		logger.Err(ctx, err).Msg("Failed to query agents from database.")
		return err
	}
	host := wasmhost.GetWasmHost(ctx)
	for _, agent := range agents {
		if !runningAgents[agent.Id] {
			go func(f_ctx context.Context, agentId string, agentName string) {
				if err := spawnActorForAgent(host, plugin, agentId, agentName, false); err != nil {
					logger.Err(f_ctx, err).Msgf("Failed to spawn actor for agent %s.", agentId)
				}
			}(ctx, agent.Id, agent.Name)
		}
	}

	return nil
}

func beforeShutdown(ctx context.Context) {
	logger.Info(ctx).Msg("Actor system shutting down...")

	// stop all agent actors before shutdown so they can suspend properly
	for _, pid := range _actorSystem.Actors() {
		if _, ok := pid.Actor().(*wasmAgentActor); ok {

			// pass the pid so it can be used during shutdown as an event sender
			ctx := context.WithValue(ctx, pidContextKey{}, pid)
			if err := pid.Shutdown(ctx); err != nil {
				logger.Err(ctx, err).Msgf("Failed to shutdown actor %s.", pid.Name())
			}
		}
	}
}

func Shutdown(ctx context.Context) {
	if _actorSystem == nil {
		return
	}

	beforeShutdown(ctx)

	if _remoting != nil {
		_remoting.Close()
		_remoting = nil
	}

	if err := _actorSystem.Stop(ctx); err != nil {
		logger.Err(ctx, err).Msg("Failed to shutdown actor system.")
	}

	logger.Info(ctx).Msg("Actor system shutdown complete.")
}
