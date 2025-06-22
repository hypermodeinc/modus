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
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	goakt "github.com/tochemey/goakt/v3/actor"
)

var _actorSystem goakt.ActorSystem

func Initialize(ctx context.Context) {

	wasmExt := &wasmExtension{
		host: wasmhost.GetWasmHost(ctx),
	}

	opts := []goakt.Option{
		goakt.WithLogger(newActorLogger(logger.Get(ctx))),
		goakt.WithCoordinatedShutdown(beforeShutdown),
		goakt.WithPubSub(),
		goakt.WithActorInitTimeout(10 * time.Second), // TODO: adjust this value, or make it configurable
		goakt.WithActorInitMaxRetries(1),             // TODO: adjust this value, or make it configurable
		goakt.WithExtensions(wasmExt),
	}
	opts = append(opts, clusterOptions(ctx)...)

	if actorSystem, err := goakt.NewActorSystem("modus", opts...); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to create actor system.")
	} else if err := actorSystem.Start(ctx); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to start actor system.")
	} else if err := actorSystem.Inject(&wasmAgentInfo{}); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to inject wasm agent info into actor system.")
	} else {
		_actorSystem = actorSystem
	}

	waitForClusterSync()

	logger.Info(ctx).Msg("Actor system started.")

	pluginmanager.RegisterPluginLoadedCallback(loadAgentActors)
}

func loadAgentActors(ctx context.Context, plugin *plugins.Plugin) error {
	// restart local actors that are already running, which will reload the plugin
	actors := _actorSystem.Actors()
	localAgents := make(map[string]bool, len(actors))
	for _, pid := range actors {
		if a, ok := pid.Actor().(*wasmAgentActor); ok {
			localAgents[a.agentId] = true
			if err := goakt.Tell(ctx, pid, &messages.RestartAgent{}); err != nil {
				logger.Err(ctx, err).Str("agent_id", a.agentId).Msg("Failed to send restart agent message to actor.")
			}
		}
	}

	// spawn actors for agents with state in the database, that are not already running
	// check both locally and on remote nodes in the cluster
	logger.Debug(ctx).Msg("Restoring agent actors from database.")
	agents, err := db.QueryActiveAgents(ctx)
	if err != nil {
		return fmt.Errorf("failed to query active agents: %w", err)
	}
	inCluster := _actorSystem.InCluster()
	for _, agent := range agents {
		if !localAgents[agent.Id] {
			if inCluster {
				actorName := getActorName(agent.Id)
				if exists, err := _actorSystem.ActorExists(ctx, actorName); err != nil {
					logger.Err(ctx, err).Msgf("Failed to check if actor %s exists in cluster.", actorName)
				} else if exists {
					// if the actor already exists in the cluster, skip spawning it
					continue
				}
			}
			if err := spawnActorForAgent(ctx, plugin.Name(), agent.Id, agent.Name, false); err != nil {
				logger.Err(ctx, err).Msgf("Failed to spawn actor for agent %s.", agent.Id)
			}
		}
	}

	return nil
}

func beforeShutdown(ctx context.Context) error {
	_actorSystem.Logger().(*actorLogger).shuttingDown = true
	logger.Info(ctx).Msg("Actor system shutting down...")
	actors := _actorSystem.Actors()

	// Suspend all local running agent actors first, which allows them to gracefully stop and persist their state.
	// In cluster mode, this will also allow the actor to resume on another node after this node shuts down.
	for _, pid := range actors {
		if actor, ok := pid.Actor().(*wasmAgentActor); ok && pid.IsRunning() {
			if actor.status == AgentStatusRunning {
				ctx := actor.augmentContext(ctx, pid)
				if err := actor.suspendAgent(ctx); err != nil {
					logger.Err(ctx, err).Str("agent_id", actor.agentId).Msg("Failed to suspend agent actor.")
				}
			}
		}
	}

	// Then shut down subscription actors.  They will have received the suspend message already.
	for _, pid := range actors {
		if _, ok := pid.Actor().(*subscriptionActor); ok && pid.IsRunning() {
			if err := pid.Shutdown(ctx); err != nil {
				logger.Err(ctx, err).Msgf("Failed to shutdown actor %s.", pid.Name())
			}
		}
	}

	waitForClusterSync()

	// then allow the actor system to continue with its shutdown process
	return nil
}

func waitForClusterSync() {
	if clusterEnabled() {
		time.Sleep(peerSyncInterval() * 2)
	}
}

func Shutdown(ctx context.Context) {
	if _actorSystem == nil {
		logger.Fatal(ctx).Msg("Actor system is not initialized, cannot shutdown.")
	}

	if err := _actorSystem.Stop(ctx); err != nil {
		logger.Err(ctx, err).Msg("Failed to shutdown actor system.")
	}

	logger.Info(ctx).Msg("Actor system shutdown complete.")
}

const wasmExtensionId = "wasm"

type wasmExtension struct {
	host wasmhost.WasmHost
}

func (w *wasmExtension) ID() string {
	return wasmExtensionId
}
