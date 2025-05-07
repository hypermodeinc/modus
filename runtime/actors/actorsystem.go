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

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/plugins"

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

	pluginmanager.RegisterPluginLoadedCallback(reloadAgentActors)
}

func reloadAgentActors(ctx context.Context, plugin *plugins.Plugin) error {
	for _, pid := range _actorSystem.Actors() {
		if actor, ok := pid.Actor().(*WasmAgentActor); ok {
			if err := actor.reloadModule(ctx, plugin); err != nil {
				return err
			}
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
