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
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	goakt "github.com/tochemey/goakt/v3/actor"
)

var _actorSystem goakt.ActorSystem

func Initialize(ctx context.Context) {

	actorLogger := newActorLogger(logger.Get(ctx))

	actorSystem, err := goakt.NewActorSystem("modus",
		goakt.WithLogger(actorLogger),
		goakt.WithCoordinatedShutdown(beforeShutdown),
		goakt.WithPassivationDisabled(),            // TODO: enable passivation. Requires a persistence store in production for agent state.
		goakt.WithActorInitTimeout(10*time.Second), // TODO: adjust this value, or make it configurable
		goakt.WithActorInitMaxRetries(1))           // TODO: adjust this value, or make it configurable

	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to create actor system.")
	}

	if err := actorSystem.Start(ctx); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to start actor system.")
	}

	_actorSystem = actorSystem

	logger.Info(ctx).Msg("Actor system started.")

	pluginmanager.RegisterPluginLoadedCallback(loadAgentActors)
}

func loadAgentActors(ctx context.Context, plugin *plugins.Plugin) error {
	// reload modules for actors that are already running
	actors := _actorSystem.Actors()
	runningAgents := make(map[string]bool, len(actors))
	for _, pid := range actors {
		if actor, ok := pid.Actor().(*wasmAgentActor); ok {
			runningAgents[actor.agentId] = true
			if err := actor.reloadModule(ctx, plugin); err != nil {
				return err
			}
		}
	}

	// spawn actors for agents with state in the database, that are not already running
	// TODO: when we scale out with GoAkt cluster mode, we'll need to decide which node is responsible for spawning the actor
	agents, err := db.QueryActiveAgents(ctx)
	if err != nil {
		logger.Err(ctx, err).Msg("Failed to query agents from database.")
		return err
	}
	host := wasmhost.GetWasmHost(ctx)
	for _, agent := range agents {
		if !runningAgents[agent.Id] {
			spawnActorForAgent(host, plugin, agent.Id, agent.Name, true, &agent.Data)
		}
	}

	return nil
}

func beforeShutdown(ctx context.Context) error {
	logger.Info(ctx).Msg("Actor system shutting down...")
	return nil
}

func Shutdown(ctx context.Context) {
	if _actorSystem == nil {
		return
	}

	if err := _actorSystem.Stop(ctx); err != nil {
		logger.Err(ctx, err).Msg("Failed to shutdown actor system.")
	}

	logger.Info(ctx).Msg("Actor system shutdown complete.")
}
