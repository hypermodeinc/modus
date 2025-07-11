/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	goakt "github.com/tochemey/goakt/v3/actor"
)

var _actorSystem goakt.ActorSystem

func Initialize(ctx context.Context) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	wasmExt := &wasmExtension{
		host: wasmhost.GetWasmHost(ctx),
	}

	opts := []goakt.Option{
		goakt.WithLogger(newActorLogger(logger.Get(ctx))),
		goakt.WithCoordinatedShutdown(&shutdownHook{}),
		goakt.WithPubSub(),
		goakt.WithActorInitTimeout(10 * time.Second), // TODO: adjust this value, or make it configurable
		goakt.WithActorInitMaxRetries(1),             // TODO: adjust this value, or make it configurable
		goakt.WithExtensions(wasmExt),
	}
	opts = append(opts, clusterOptions(ctx)...)

	actorSystem, err := goakt.NewActorSystem("modus", opts...)
	if err != nil {
		const msg = "Failed to create actor system."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Fatal(ctx, err).Msg(msg)
	}

	if err := startActorSystem(ctx, actorSystem); err != nil {
		const msg = "Failed to start actor system."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Fatal(ctx, err).Msg(msg)
	}

	if err := actorSystem.Inject(&wasmAgentInfo{}); err != nil {
		const msg = "Failed to inject wasm agent info into actor system."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Fatal(ctx, err).Msg(msg)
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
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	// restart local actors that are already running, which will reload the plugin
	actors := _actorSystem.Actors()
	for _, pid := range actors {
		if a, ok := pid.Actor().(*wasmAgentActor); ok {
			if err := goakt.Tell(ctx, pid, &messages.RestartAgent{}); err != nil {
				const msg = "Failed to send restart agent message to actor."
				sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("agent_id", a.agentId))
				logger.Error(ctx, err).Str("agent_id", a.agentId).Msg(msg)
			}
		}
	}

	// do this in a goroutine to avoid blocking the cluster engine startup
	go func() {
		if err := restoreAgentActors(ctx, plugin.Name()); err != nil {
			const msg = "Failed to restore agent actors."
			sentryutils.CaptureError(ctx, err, msg)
			logger.Error(ctx, err).Msg(msg)
		}
	}()

	return nil
}

// restoreAgentActors spawn actors for agents with state in the database, that are not already running
func restoreAgentActors(ctx context.Context, pluginName string) error {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
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
			const msg = "Failed to check if agent actor exists."
			sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("agent_id", agent.Id))
			logger.Error(ctx, err).Str("agent_id", agent.Id).Msg(msg)
		} else if !exists {
			err := spawnActorForAgent(ctx, pluginName, agent.Id, agent.Name, false)
			if err != nil {
				const msg = "Failed to spawn actor for agent."
				sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("agent_id", agent.Id))
				logger.Error(ctx, err).Str("agent_id", agent.Id).Msg(msg)
			}
		}
	}

	return nil
}

// Waits for the peer sync interval to pass, allowing time for the actor system to synchronize its
// list of actors with the remote nodes in the cluster. Cancels early if the context is done.
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
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	if _actorSystem == nil {
		const msg = "Actor system is not initialized, cannot shutdown."
		sentryutils.CaptureError(ctx, nil, msg)
		logger.Fatal(ctx).Msg(msg)
	}

	if err := _actorSystem.Stop(ctx); err != nil {
		const msg = "Failed to shutdown actor system."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Error(ctx, err).Msg(msg)
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

type shutdownHook struct {
}

func (sh *shutdownHook) Execute(ctx context.Context, actorSystem goakt.ActorSystem) error {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	logger.Info(ctx).Msg("Actor system shutting down...")
	actors := actorSystem.Actors()

	// Suspend all local running agent actors first, which allows them to gracefully stop and persist their state.
	// In cluster mode, this will also allow the actor to resume on another node after this node shuts down.
	for _, pid := range actors {
		if actor, ok := pid.Actor().(*wasmAgentActor); ok && pid.IsRunning() {
			if actor.status == AgentStatusRunning {
				ctx := actor.augmentContext(ctx, pid)
				if err := actor.suspendAgent(ctx); err != nil {
					const msg = "Failed to suspend agent actor."
					sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("agent_id", actor.agentId))
					logger.Error(ctx, err).Str("agent_id", actor.agentId).Msg(msg)
				}
			}
		}
	}

	// Then shut down subscription actors. They will have received the suspend message already.
	for _, pid := range actors {
		if a, ok := pid.Actor().(*subscriptionActor); ok && pid.IsRunning() {
			if err := pid.Shutdown(ctx); err != nil {
				const msg = "Failed to shut down subscription actor."
				sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("agent_id", a.agentId))
				logger.Error(ctx, err).Str("agent_id", a.agentId).Msg(msg)
			}
		}
	}

	// Then allow the actor system to continue with its shutdown process.
	return nil
}

func (sh *shutdownHook) Recovery() *goakt.ShutdownHookRecovery {
	return goakt.NewShutdownHookRecovery(
		goakt.WithShutdownHookRetry(2, 2*time.Second),
		goakt.WithShutdownHookRecoveryStrategy(goakt.ShouldRetryAndSkip),
	)
}
