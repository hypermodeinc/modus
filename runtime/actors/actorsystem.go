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
	"os"
	"strconv"
	"time"

	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	goakt "github.com/tochemey/goakt/v3/actor"
	goakt_static "github.com/tochemey/goakt/v3/discovery/static"
	goakt_remote "github.com/tochemey/goakt/v3/remote"
	"github.com/travisjeffery/go-dynaport"
)

var _actorSystem goakt.ActorSystem

func Initialize(ctx context.Context) {

	opts := []goakt.Option{
		goakt.WithLogger(newActorLogger(logger.Get(ctx))),
		goakt.WithCoordinatedShutdown(beforeShutdown),
		goakt.WithPubSub(),
		goakt.WithActorInitTimeout(10 * time.Second), // TODO: adjust this value, or make it configurable
		goakt.WithActorInitMaxRetries(1),             // TODO: adjust this value, or make it configurable

		// for now, keep passivation disabled so that agents can perform long-running tasks without the actor stopping.
		// TODO: figure out how to deal with this better
		goakt.WithPassivationDisabled(),
	}

	// NOTE: we're not relying on cluster mode yet.  The below code block is for future use and testing purposes only.
	if clusterMode, _ := strconv.ParseBool(os.Getenv("MODUS_USE_CLUSTER_MODE")); clusterMode {
		// TODO: static discovery should really only be used for local development and testing.
		// In production, we should use a more robust discovery mechanism, such as Kubernetes or NATS.
		// See https://tochemey.gitbook.io/goakt/features/service-discovery

		// We just get three random ports for now.
		// In prod, these will need to be configured so they are consistent across all nodes.
		ports := dynaport.Get(3)
		var gossip_port = ports[0]
		var peers_port = ports[1]
		var remoting_port = ports[2]

		disco := goakt_static.NewDiscovery(&goakt_static.Config{
			Hosts: []string{
				fmt.Sprintf("localhost:%d", gossip_port),
			},
		})

		opts = append(opts,
			goakt.WithRemote(goakt_remote.NewConfig("localhost", remoting_port)),
			goakt.WithCluster(goakt.NewClusterConfig().
				WithDiscovery(disco).
				WithDiscoveryPort(gossip_port).
				WithPeersPort(peers_port).
				WithKinds(&wasmAgentActor{}, &subscriptionActor{}),
			),
		)
	}

	if actorSystem, err := goakt.NewActorSystem("modus", opts...); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to create actor system.")
	} else if err := actorSystem.Start(ctx); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to start actor system.")
	} else {
		_actorSystem = actorSystem
	}

	logger.Info(ctx).Msg("Actor system started.")

	pluginmanager.RegisterPluginLoadedCallback(loadAgentActors)
}

func loadAgentActors(ctx context.Context, plugin *plugins.Plugin) error {
	// restart actors that are already running, giving them the new plugin instance
	actors := _actorSystem.Actors()
	runningAgents := make(map[string]bool, len(actors))
	for _, pid := range actors {
		if actor, ok := pid.Actor().(*wasmAgentActor); ok {
			runningAgents[actor.agentId] = true
			actor.plugin = plugin
			if err := pid.Restart(ctx); err != nil {
				logger.Err(ctx, err).Msgf("Failed to restart actor for agent %s.", actor.agentId)
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
			spawnActorForAgentAsync(host, plugin, agent.Id, agent.Name, false)
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
