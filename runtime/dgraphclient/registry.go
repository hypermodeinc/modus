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
	"sync"

	"github.com/hypermodeinc/modus/pkg/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
)

var dgr = newDgraphRegistry()

type dgraphRegistry struct {
	sync.RWMutex
	dgraphConnectorCache map[string]*dgraphConnector
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
		dgraphConnectorCache: make(map[string]*dgraphConnector),
	}
}

func ShutdownConns() {
	dgr.Lock()
	defer dgr.Unlock()
	for _, ds := range dgr.dgraphConnectorCache {
		ds.conn.Close()
	}
	clear(dgr.dgraphConnectorCache)
}

func (dr *dgraphRegistry) getDgraphConnector(ctx context.Context, dgName string) (*dgraphConnector, error) {
	dr.RLock()
	ds, ok := dr.dgraphConnectorCache[dgName]
	dr.RUnlock()
	if ok {
		return ds, nil
	}

	dr.Lock()
	defer dr.Unlock()

	if ds, ok := dr.dgraphConnectorCache[dgName]; ok {
		return ds, nil
	}

	for name, info := range manifestdata.GetManifest().Hosts {
		if name != dgName {
			continue
		}

		if info.HostType() != manifest.HostTypeDgraph {
			return nil, fmt.Errorf("host %s is not a dgraph host", dgName)
		}

		host := info.(manifest.DgraphHostInfo)
		if host.GrpcTarget == "" {
			return nil, fmt.Errorf("dgraph host %s has empty GrpcTarget", dgName)
		}

		var opts []grpc.DialOption

		if host.Key != "" {
			hostKey, err := secrets.ApplyHostSecretsToString(ctx, info, host.Key)
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
				grpc.WithPerRPCCredentials(&authCreds{hostKey}),
			}
		} else if strings.Split(host.GrpcTarget, ":")[0] != "localhost" {
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

		conn, err := grpc.NewClient(host.GrpcTarget, opts...)
		if err != nil {
			return nil, err
		}

		ds := &dgraphConnector{
			conn:     conn,
			dgClient: dgo.NewDgraphClient(api.NewDgraphClient(conn)),
		}
		dr.dgraphConnectorCache[dgName] = ds
		return ds, nil
	}

	return nil, fmt.Errorf("dgraph host %s not found", dgName)
}
