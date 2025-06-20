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
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"

	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/discovery"
	"github.com/tochemey/goakt/v3/discovery/kubernetes"
	"github.com/tochemey/goakt/v3/discovery/nats"
	"github.com/tochemey/goakt/v3/remote"
	"github.com/travisjeffery/go-dynaport"
)

func clusterOptions(ctx context.Context) []goakt.Option {

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

	var disco discovery.Provider
	switch clusterMode {
	case clusterModeNats:
		natsUrl := clusterNatsUrl()
		clusterHost := clusterHost()
		logger.Info(ctx).
			Str("cluster_host", clusterHost).
			Str("nats_server", natsUrl).
			Msg("Using NATS for node discovery.")

		disco = nats.NewDiscovery(&nats.Config{
			NatsSubject:   "modus-gossip",
			NatsServer:    natsUrl,
			Host:          clusterHost,
			DiscoveryPort: discoveryPort,
		})

	case clusterModeKubernetes:
		namespace, ok := app.KubernetesNamespace()
		if !ok {
			logger.Fatal(ctx).
				Msg("Kubernetes cluster mode enabled, but a Kubernetes namespace was not found. Ensure running in a Kubernetes environment.")
			return nil
		}

		logger.Info(ctx).
			Str("namespace", namespace).
			Msg("Using Kubernetes for node discovery.")

		disco = kubernetes.NewDiscovery(&kubernetes.Config{
			Namespace:         namespace,
			PodLabels:         getPodLabels(),
			DiscoveryPortName: "discovery-port",
			RemotingPortName:  "remoting-port",
			PeersPortName:     "peers-port",
		})

	default:
		panic("Unsupported cluster mode: " + clusterMode.String())
	}

	var remotingHost string
	if app.IsDevEnvironment() {
		// only bind to localhost in development
		remotingHost = "127.0.0.1"
	} else {
		// otherwise bind to all interfaces
		remotingHost = "0.0.0.0"
	}

	return []goakt.Option{
		goakt.WithRemote(remote.NewConfig(remotingHost, remotingPort)),
		goakt.WithCluster(goakt.NewClusterConfig().
			WithDiscovery(disco).
			WithDiscoveryPort(discoveryPort).
			WithPeersPort(peersPort).
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
		return "localhost"
	} else {
		// this hack gets the same IP that the remoting system would bind to by default
		rc := remote.NewConfig("0.0.0.0", 0)
		_ = rc.Sanitize()
		return rc.BindAddr()
	}
}

func clusterPorts() (discoveryPort, remotingPort, peersPort int) {

	// Get default ports dynamically
	ports := dynaport.Get(3)
	discoveryPort = ports[0]
	remotingPort = ports[1]
	peersPort = ports[2]

	// Override with environment variables if set
	discoveryPort = getPortFromEnv("MODUS_CLUSTER_DISCOVERY_PORT", discoveryPort)
	remotingPort = getPortFromEnv("MODUS_CLUSTER_REMOTING_PORT", remotingPort)
	peersPort = getPortFromEnv("MODUS_CLUSTER_PEERS_PORT", peersPort)

	return
}

func getPortFromEnv(envVar string, defaultPort int) int {
	portStr := os.Getenv(envVar)
	if portStr == "" {
		return defaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		logger.Warnf("Invalid value for %s. Using %d instead.", envVar, defaultPort)
		return defaultPort
	}

	return port
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
