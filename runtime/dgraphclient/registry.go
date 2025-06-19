/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraphclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"

	"github.com/dgraph-io/dgo/v250"
	"github.com/puzpuzpuz/xsync/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var dgr = newDgraphRegistry()

type dgraphRegistry struct {
	cache *xsync.Map[string, *dgraphConnector]
}

func newDgraphRegistry() *dgraphRegistry {
	return &dgraphRegistry{
		cache: xsync.NewMap[string, *dgraphConnector](),
	}
}

func ShutdownConns() {
	dgr.cache.Range(func(key string, _ *dgraphConnector) bool {
		if connector, ok := dgr.cache.LoadAndDelete(key); ok {
			connector.dgClient.Close()
		}
		return true
	})
}

func (dr *dgraphRegistry) getDgraphConnector(ctx context.Context, dgName string) (*dgraphConnector, error) {
	var creationErr error
	ds, _ := dr.cache.LoadOrCompute(dgName, func() (*dgraphConnector, bool) {
		conn, err := createConnector(ctx, dgName)
		if err != nil {
			creationErr = err
			return nil, true
		}
		return conn, false
	})
	return ds, creationErr
}

func createConnector(ctx context.Context, dgName string) (*dgraphConnector, error) {
	man := manifestdata.GetManifest()
	info, ok := man.Connections[dgName]
	if !ok {
		return nil, fmt.Errorf("dgraph connection [%s] not found", dgName)
	}

	if info.ConnectionType() != manifest.ConnectionTypeDgraph {
		return nil, fmt.Errorf("[%s] is not a dgraph connection", dgName)
	}

	connection := info.(manifest.DgraphConnectionInfo)
	if connection.ConnStr != "" {
		if connection.GrpcTarget != "" {
			return nil, fmt.Errorf("dgraph connection [%s] has both connString and grpcTarget (use one or the other)", dgName)
		} else if connection.Key != "" {
			return nil, fmt.Errorf("dgraph connection [%s] has both connString and key (the key should be part of the connection string)", dgName)
		}
		connStr, err := secrets.ApplySecretsToString(ctx, info, connection.ConnStr)
		if err != nil {
			return nil, err
		}
		return connectWithConnectionString(connStr)
	} else if connection.GrpcTarget != "" {
		target, err := secrets.ApplySecretsToString(ctx, info, connection.GrpcTarget)
		if err != nil {
			return nil, err
		}
		key, err := secrets.ApplySecretsToString(ctx, info, connection.Key)
		if err != nil {
			return nil, err
		}
		return connectWithGrpcTarget(target, key)
	} else {
		return nil, fmt.Errorf("dgraph connection [%s] needs either a connString or a grpcTarget", dgName)
	}
}

func connectWithConnectionString(connStr string) (*dgraphConnector, error) {
	if dgClient, err := dgo.Open(connStr); err != nil {
		return nil, err
	} else {
		return newDgraphConnector(dgClient), nil
	}
}

func connectWithGrpcTarget(target string, key string) (*dgraphConnector, error) {
	var opts []dgo.ClientOption
	if strings.Split(target, ":")[0] == "localhost" {
		opts = []dgo.ClientOption{
			dgo.WithGrpcOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		}
	} else if key == "" {
		opts = []dgo.ClientOption{
			dgo.WithSystemCertPool(),
		}
	} else if strings.Contains(strings.ToLower(target), "cloud.dgraph.io") {
		opts = []dgo.ClientOption{
			dgo.WithSystemCertPool(),
			dgo.WithDgraphAPIKey(key),
		}
	} else {
		opts = []dgo.ClientOption{
			dgo.WithSystemCertPool(),
			dgo.WithBearerToken(key),
		}
	}

	if dgClient, err := dgo.NewClient(target, opts...); err != nil {
		return nil, err
	} else {
		return newDgraphConnector(dgClient), nil
	}
}
