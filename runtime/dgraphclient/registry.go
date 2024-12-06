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
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/puzpuzpuz/xsync/v3"
)

var dgr = newDgraphRegistry()

type dgraphRegistry struct {
	cache *xsync.MapOf[string, *dgraphConnector]
}

type authCreds struct {
	token string
}

func (a *authCreds) GetRequestMetadata(ctx context.Context, uri ...string) (
	map[string]string, error) {

	return map[string]string{"Authorization": a.token}, nil
}

func (a *authCreds) RequireTransportSecurity() bool {
	return true
}

func newDgraphRegistry() *dgraphRegistry {
	return &dgraphRegistry{
		cache: xsync.NewMapOf[string, *dgraphConnector](),
	}
}

func ShutdownConns() {
	dgr.cache.Range(func(key string, _ *dgraphConnector) bool {
		if ds, ok := dgr.cache.LoadAndDelete(key); ok {
			ds.conn.Close()
		}
		return true
	})
}

func (dr *dgraphRegistry) getDgraphConnector(ctx context.Context, dgName string) (*dgraphConnector, error) {
	var creationErr error
	ds, _ := dr.cache.LoadOrCompute(dgName, func() *dgraphConnector {
		conn, err := createConnector(ctx, dgName)
		if err != nil {
			creationErr = err
			return nil
		}
		return conn
	})

	if creationErr != nil {
		dr.cache.Delete(dgName)
		return nil, creationErr
	}

	return ds, nil
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
	if connection.GrpcTarget == "" {
		return nil, fmt.Errorf("dgraph connection [%s] has empty GrpcTarget", dgName)
	}

	var opts []grpc.DialOption

	if connection.Key != "" {
		conKey, err := secrets.ApplySecretsToString(ctx, info, connection.Key)
		if err != nil {
			return nil, err
		}

		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(&authCreds{conKey}),
		}
	} else if strings.Split(connection.GrpcTarget, ":")[0] != "localhost" {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}
	} else {
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	}

	conn, err := grpc.NewClient(connection.GrpcTarget, opts...)
	if err != nil {
		return nil, err
	}

	ds := &dgraphConnector{
		conn:     conn,
		dgClient: dgo.NewDgraphClient(api.NewDgraphClient(conn)),
	}

	return ds, nil
}
