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
	"net/url"
	"os"
	"strconv"

	"github.com/hypermodeinc/modus/runtime/logger"

	goakt "github.com/tochemey/goakt/v3/actor"
	goakt_nats "github.com/tochemey/goakt/v3/discovery/nats"
	goakt_remote "github.com/tochemey/goakt/v3/remote"
	"github.com/travisjeffery/go-dynaport"
)

func clusterOptions() []goakt.Option {
	if !clusterModeEnabled() {
		return nil
	}

	ports := dynaport.Get(3)
	discoveryPort := ports[0]
	remotingPort := ports[1]
	peersPort := ports[2]

	disco := goakt_nats.NewDiscovery(&goakt_nats.Config{
		ActorSystemName: actorSystemName,
		ApplicationName: "Modus Runtime",
		NatsSubject:     "modus-gossip",
		NatsServer:      clusterNatsUrl(),
		Host:            "localhost",
		DiscoveryPort:   discoveryPort,
		// Timeout:         0,
		// MaxJoinAttempts: 0,
		// ReconnectWait:   0,
	})

	return []goakt.Option{
		goakt.WithRemote(goakt_remote.NewConfig("localhost", remotingPort)),
		goakt.WithCluster(goakt.NewClusterConfig().
			WithDiscovery(disco).
			WithDiscoveryPort(discoveryPort).
			WithPeersPort(peersPort).
			WithKinds(&wasmAgentActor{}, &subscriptionActor{}),
		),
	}

}

func clusterModeEnabled() bool {
	ok, _ := strconv.ParseBool(os.Getenv("MODUS_USE_CLUSTER_MODE"))
	return ok
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
