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
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/discovery"
	"github.com/tochemey/goakt/v3/discovery/kubernetes"
	"github.com/tochemey/goakt/v3/discovery/nats"
	"github.com/tochemey/goakt/v3/remote"
	"github.com/travisjeffery/go-dynaport"
)

func clusterOptions(ctx context.Context) []goakt.Option {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	clusterMode := clusterMode()
	if clusterMode == clusterModeNone {
		if !app.IsDevEnvironment() {
			logger.Warnf("Cluster mode is disabled, which is not recommended for production environments. Set MODUS_CLUSTER_MODE to enable clustering.")
		}
		return nil
	}

	discoveryPort, remotingPort, peersPort := clusterPorts()
	logger.Info(ctx).
		Str("cluster_mode", clusterMode.String()).
		Int("discovery_port", discoveryPort).
		Int("remoting_port", remotingPort).
		Int("peers_port", peersPort).
		Msg("Clustering enabled.")

	disco, err := newDiscoveryProvider(ctx, clusterMode, discoveryPort)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to create cluster discovery provider.")
	}

	return []goakt.Option{
		goakt.WithRemote(remote.NewConfig(remotingHost(), remotingPort)),
		goakt.WithCluster(goakt.NewClusterConfig().
			WithDiscovery(disco).
			WithDiscoveryPort(discoveryPort).
			WithPeersPort(peersPort).
			WithReadTimeout(readTimeout()).
			WithWriteTimeout(writeTimeout()).
			WithPartitionCount(partitionCount()).
			WithClusterStateSyncInterval(nodesSyncInterval()).
			WithPeersStateSyncInterval(peerSyncInterval()).
			WithKinds(&wasmAgentActor{}, &subscriptionActor{}),
		),
	}
}

type goaktClusterMode int

const (
	clusterModeNone goaktClusterMode = iota
	clusterModeNats
	clusterModeKubernetes
)

func (c goaktClusterMode) String() string {
	switch c {
	case clusterModeNone:
		return "none"
	case clusterModeNats:
		return "NATS"
	case clusterModeKubernetes:
		return "Kubernetes"
	default:
		return "unknown"
	}
}

func parseClusterMode(mode string) goaktClusterMode {
	switch strings.ToLower(mode) {
	case "none", "":
		return clusterModeNone
	case "nats":
		return clusterModeNats
	case "kubernetes", "k8s":
		return clusterModeKubernetes
	default:
		logger.Warnf("Unknown cluster mode: '%s'. Defaulting to 'none'.", mode)
		return clusterModeNone
	}
}

func clusterMode() goaktClusterMode {
	return parseClusterMode(os.Getenv("MODUS_CLUSTER_MODE"))
}

func clusterEnabled() bool {
	return clusterMode() != clusterModeNone
}

func clusterNatsUrl() string {
	const envVar = "MODUS_CLUSTER_NATS_URL"
	const defaultNatsUrl = "nats://localhost:4222"
	urlStr := os.Getenv(envVar)
	if urlStr == "" {
		logger.Warnf("%s not set. Using default: %s", envVar, defaultNatsUrl)
		return defaultNatsUrl
	}
	if _, err := url.Parse(urlStr); err != nil {
		logger.Warnf("Invalid URL for %s. Using default: %s", envVar, defaultNatsUrl)
		return defaultNatsUrl
	}

	return urlStr
}

func clusterHost() string {
	const envVar = "MODUS_CLUSTER_HOST"
	if host := os.Getenv(envVar); host != "" {
		if _, err := url.Parse("http://" + host); err != nil {
			logger.Fatalf("Invalid value for %s: %s.", envVar, host)
		}
		return host
	}

	if app.IsDevEnvironment() {
		// Note, forcing IPv4 here avoids memberlist attempting to bind to IPv6 that we're not listening on.
		return "127.0.0.1"
	} else {
		// this hack gets the same IP that the remoting system would bind to by default
		rc := remote.NewConfig("0.0.0.0", 0)
		_ = rc.Sanitize()
		return rc.BindAddr()
	}
}

// remotingHost returns the host address to bind the remoting system to.
func remotingHost() string {
	// only bind to localhost in development
	if app.IsDevEnvironment() {
		return "127.0.0.1"
	}

	// otherwise bind to all interfaces
	return "0.0.0.0"
}

// clusterPorts returns the ports used for discovery, remoting, and peer communication in the cluster.
func clusterPorts() (discoveryPort, remotingPort, peersPort int) {

	// Get default ports dynamically, but use environment variables if set
	ports := dynaport.Get(3)
	discoveryPort = getIntFromEnv("MODUS_CLUSTER_DISCOVERY_PORT", ports[0])
	remotingPort = getIntFromEnv("MODUS_CLUSTER_REMOTING_PORT", ports[1])
	peersPort = getIntFromEnv("MODUS_CLUSTER_PEERS_PORT", ports[2])

	return
}

// peerSyncInterval returns the interval at which the actor system will sync its list of actors to other nodes across the cluster.
// We use a tight sync interval of 1 second by default, to ensure quick peer discovery as agents are added or removed.
//
// This value is also used for a sleep both on system startup and when spawning a new agent actor,
// so it needs to be low enough to not be noticed by the user.
func peerSyncInterval() time.Duration {
	return getDurationFromEnv("MODUS_CLUSTER_PEER_SYNC_SECONDS", 1, time.Second)
}

// nodesSyncInterval returns the interval at which the cluster forces a resync of the list of active nodes across the cluster.
// This matters only with regard to nodes going down unexpectedly, as other nodes in the cluster will not be aware of the change until the next sync.
// It does not affect anything if a node is gracefully shut down, as that will be communicated immediately during the shutdown process.
//
// On each interval, the node will sync its list of nodes with the cluster, and update its local state accordingly.
// The default is 10 seconds, which is a reasonable balance between responsiveness and network overhead.
func nodesSyncInterval() time.Duration {
	return getDurationFromEnv("MODUS_CLUSTER_NODES_SYNC_SECONDS", 10, time.Second)
}

// partitionCount returns the number of partitions the cluster will use for actor distribution.
// It must be a prime number to work properly with the actor system's hashing algorithm.
// It must be greater than the number of nodes in the cluster, but not too large to avoid excessive overhead.
// In testing, 23 is the highest that works well with the other default timing constraints.
// We'll use a slightly lower default of 13, which is still a prime number and should work well for most clusters.
// The GoAkt default is 271, but this has been found to lead to other errors in practice.
func partitionCount() uint64 {
	return uint64(getIntFromEnv("MODUS_CLUSTER_PARTITION_COUNT", 13))
}

// readTimeout returns the duration to wait for a cluster read operation before timing out.
// The default is 1 second, which should usually not need to be changed.
func readTimeout() time.Duration {
	return getDurationFromEnv("MODUS_CLUSTER_READ_TIMEOUT_SECONDS", 1, time.Second)
}

// writeTimeout returns the duration to wait for a cluster write operation before timing out.
// The default is 1 second, which should usually not need to be changed.
func writeTimeout() time.Duration {
	return getDurationFromEnv("MODUS_CLUSTER_WRITE_TIMEOUT_SECONDS", 1, time.Second)
}

func getPodLabels() map[string]string {
	// example value: "app.kubernetes.io/name=modus,app.kubernetes.io/component=runtime"
	if labels := os.Getenv("MODUS_CLUSTER_POD_LABELS"); labels != "" {
		podLabels := make(map[string]string)
		for label := range strings.SplitSeq(labels, ",") {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				podLabels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			} else {
				logger.Warnf("Invalid pod label format: '%s'. Expected 'key=value'.", label)
			}
		}
		return podLabels
	}

	// defaults
	return map[string]string{
		"app.kubernetes.io/name":      "modus",
		"app.kubernetes.io/component": "runtime",
	}
}

func newDiscoveryProvider(ctx context.Context, clusterMode goaktClusterMode, discoveryPort int) (discovery.Provider, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	switch clusterMode {
	case clusterModeNats:
		natsUrl := clusterNatsUrl()
		clusterHost := clusterHost()

		logger.Info(ctx).
			Str("cluster_host", clusterHost).
			Str("nats_server", natsUrl).
			Msg("Using NATS for node discovery.")

		disco := nats.NewDiscovery(&nats.Config{
			NatsSubject:   "modus-gossip",
			NatsServer:    natsUrl,
			Host:          clusterHost,
			DiscoveryPort: discoveryPort,
		})

		return wrapProvider(ctx, disco), nil

	case clusterModeKubernetes:
		namespace, ok := app.KubernetesNamespace()
		if !ok {
			return nil, errors.New("Kubernetes cluster mode enabled, but namespace was not found")
		}

		logger.Info(ctx).
			Str("namespace", namespace).
			Msg("Using Kubernetes for node discovery.")

		disco := kubernetes.NewDiscovery(&kubernetes.Config{
			Namespace:         namespace,
			PodLabels:         getPodLabels(),
			DiscoveryPortName: "discovery-port",
			RemotingPortName:  "remoting-port",
			PeersPortName:     "peers-port",
		})

		return wrapProvider(ctx, disco), nil
	}

	return nil, fmt.Errorf("unsupported cluster mode: %s", clusterMode)
}

// wrapProvider wraps a discovery provider to add Sentry tracing to its methods.
func wrapProvider(ctx context.Context, provider discovery.Provider) discovery.Provider {
	if provider == nil {
		return nil
	}

	return &providerWrapper{ctx, provider}
}

// providerWrapper is a wrapper around a discovery provider that adds Sentry tracing to its methods.
type providerWrapper struct {
	ctx      context.Context
	provider discovery.Provider
}

func (w *providerWrapper) Close() error {
	span, _ := utils.NewSentrySpanForCurrentFunc(w.ctx)
	defer span.Finish()

	return w.provider.Close()
}

func (w *providerWrapper) Deregister() error {
	span, _ := utils.NewSentrySpanForCurrentFunc(w.ctx)
	defer span.Finish()

	return w.provider.Deregister()
}

func (w *providerWrapper) DiscoverPeers() ([]string, error) {
	span, _ := utils.NewSentrySpanForCurrentFunc(w.ctx)
	defer span.Finish()

	return w.provider.DiscoverPeers()
}

func (w *providerWrapper) ID() string {
	return w.provider.ID()
}

func (w *providerWrapper) Initialize() error {
	span, _ := utils.NewSentrySpanForCurrentFunc(w.ctx)
	defer span.Finish()

	return w.provider.Initialize()
}

func (w *providerWrapper) Register() error {
	span, _ := utils.NewSentrySpanForCurrentFunc(w.ctx)
	defer span.Finish()

	return w.provider.Register()
}
