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
	"math/rand/v2"
	"time"

	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/messages"
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	goakt "github.com/tochemey/goakt/v3/actor"
)

var _actorSystem goakt.ActorSystem

func Initialize(ctx context.Context) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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

	actorSystem, err := goakt.NewActorSystem("modus", opts...)
	if err != nil {
		logger.Fatal(ctx, err).Msg("Failed to create actor system.")
	}

	if err := startActorSystem(ctx, actorSystem); err != nil {
		logger.Fatal(ctx, err).Msg("Failed to start actor system.")
	}

	if err := actorSystem.Inject(&wasmAgentInfo{}); err != nil {
		logger.Fatal(ctx, err).Msg("Failed to inject wasm agent info into actor system.")
	}

	_actorSystem = actorSystem

	logger.Info(ctx).Msg("Actor system started.")

	pluginmanager.RegisterPluginLoadedCallback(loadAgentActors)
}

func startActorSystem(ctx context.Context, actorSystem goakt.ActorSystem) error {
	maxRetries := utils.GetIntFromEnv("MODUS_ACTOR_SYSTEM_START_MAX_RETRIES", 5)
	retryInterval := utils.GetDurationFromEnv("MODUS_ACTOR_SYSTEM_START_RETRY_INTERVAL_SECONDS", 2, time.Second)

	for i := range maxRetries {
		if err := actorSystem.Start(ctx); err != nil {
			logger.Warn(ctx, err).Int("attempt", i+1).Msgf("Failed to start actor system, retrying in %s...", retryInterval)
			time.Sleep(retryInterval)
			retryInterval *= 2 // Exponential backoff
			continue
		}

		// important: wait for the actor system to sync with the cluster before proceeding
		waitForClusterSync(ctx)

		return nil
	}

	return fmt.Errorf("failed to start actor system after %d retries", maxRetries)
}

func loadAgentActors(ctx context.Context, plugin *plugins.Plugin) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	// restart local actors that are already running, which will reload the plugin
	actors := _actorSystem.Actors()
	for _, pid := range actors {
		if a, ok := pid.Actor().(*wasmAgentActor); ok {
			if err := goakt.Tell(ctx, pid, &messages.RestartAgent{}); err != nil {
				logger.Error(ctx, err).Str("agent_id", a.agentId).Msg("Failed to send restart agent message to actor.")
			}
		}
	}

	// do this in a goroutine to avoid blocking the cluster engine startup
	go func() {
		if err := restoreAgentActors(ctx, plugin.Name()); err != nil {
			logger.Error(ctx, err).Msg("Failed to restore agent actors.")
		}
	}()

	return nil
}

// restoreAgentActors spawn actors for agents with state in the database, that are not already running
func restoreAgentActors(ctx context.Context, pluginName string) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	logger.Debug(ctx).Msg("Restoring agent actors from database.")

	// query the database for active agents
	agents, err := db.QueryActiveAgents(ctx)
	if err != nil {
		return fmt.Errorf("failed to query active agents from database: %w", err)
	}

	// shuffle the agents to help distribute the load across the cluster when multiple nodes are starting simultaneously
	rand.Shuffle(len(agents), func(i, j int) {
		agents[i], agents[j] = agents[j], agents[i]
	})

	// spawn actors for each agent that is not already running
	for _, agent := range agents {
		actorName := getActorName(agent.Id)
		if exists, err := _actorSystem.ActorExists(ctx, actorName); err != nil {
			logger.Error(ctx, err).Msgf("Failed to check if actor %s exists.", actorName)
		} else if !exists {
			err := spawnActorForAgent(ctx, pluginName, agent.Id, agent.Name, false)
			if err != nil {
				logger.Error(ctx, err).Msgf("Failed to spawn actor for agent %s.", agent.Id)
			}
		}
	}

	return nil
}

func beforeShutdown(ctx context.Context) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	logger.Info(ctx).Msg("Actor system shutting down...")
	actors := _actorSystem.Actors()

	// Suspend all local running agent actors first, which allows them to gracefully stop and persist their state.
	// In cluster mode, this will also allow the actor to resume on another node after this node shuts down.
	for _, pid := range actors {
		if actor, ok := pid.Actor().(*wasmAgentActor); ok && pid.IsRunning() {
			if actor.status == AgentStatusRunning {
				ctx := actor.augmentContext(ctx, pid)
				if err := actor.suspendAgent(ctx); err != nil {
					logger.Error(ctx, err).Str("agent_id", actor.agentId).Msg("Failed to suspend agent actor.")
				}
			}
		}
	}

	// Then shut down subscription actors.  They will have received the suspend message already.
	for _, pid := range actors {
		if _, ok := pid.Actor().(*subscriptionActor); ok && pid.IsRunning() {
			if err := pid.Shutdown(ctx); err != nil {
				logger.Error(ctx, err).Msgf("Failed to shutdown actor %s.", pid.Name())
			}
		}
	}

	// waitForClusterSync()

	// then allow the actor system to continue with its shutdown process
	return nil
}

// Waits for the peer sync interval to pass, allowing time for the actor system to synchronize its
// list of actors with the remote nodes in the cluster.  Cancels early if the context is done.
func waitForClusterSync(ctx context.Context) {
	if clusterEnabled() {
		select {
		case <-time.After(peerSyncInterval()):
		case <-ctx.Done():
			logger.Warn(context.WithoutCancel(ctx)).Msg("Context cancelled while waiting for cluster sync.")
		}
	}
}

func Shutdown(ctx context.Context) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	if _actorSystem == nil {
		logger.Fatal(ctx).Msg("Actor system is not initialized, cannot shutdown.")
	}

	if err := _actorSystem.Stop(ctx); err != nil {
		logger.Error(ctx, err).Msg("Failed to shutdown actor system.")
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
