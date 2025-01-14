/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sqlclient

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ShutdownPGPools shuts down all the PostgreSQL connection pools.
func ShutdownPGPools() {
	dsr.cache.Range(func(key string, _ *postgresqlDS) bool {
		if ds, ok := dsr.cache.LoadAndDelete(key); ok {
			ds.pool.Close()
		}
		return true
	})
}

func (r *dsRegistry) getPostgresDS(ctx context.Context, dsName string) (*postgresqlDS, error) {
	var creationErr error
	ds, _ := r.cache.LoadOrCompute(dsName, func() *postgresqlDS {
		ds, err := createDS(ctx, dsName)
		if err != nil {
			creationErr = err
			return nil
		}
		return ds
	})

	if creationErr != nil {
		r.cache.Delete(dsName)
		return nil, creationErr
	}

	return ds, nil
}

func createDS(ctx context.Context, dsName string) (*postgresqlDS, error) {
	man := manifestdata.GetManifest()
	info, ok := man.Connections[dsName]
	if !ok {
		return nil, fmt.Errorf("postgresql connection [%s] not found", dsName)
	}

	if info.ConnectionType() != manifest.ConnectionTypePostgresql {
		return nil, fmt.Errorf("[%s] is not a postgresql connection", dsName)
	}

	conf := info.(manifest.PostgresqlConnectionInfo)
	if conf.ConnStr == "" {
		return nil, fmt.Errorf("postgresql connection [%s] has empty connString", dsName)
	}

	connStr, err := secrets.ApplySecretsToString(ctx, info, conf.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to apply secrets to connection string for connection [%s]: %w", dsName, err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres connection [%s]: %w", dsName, err)
	}

	return &postgresqlDS{pool}, nil
}
