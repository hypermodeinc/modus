/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
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
	dsr.Lock()
	defer dsr.Unlock()

	for _, ds := range dsr.pgCache {
		ds.pool.Close()
	}

	clear(dsr.pgCache)
}

func (r *dsRegistry) getPGPool(ctx context.Context, dsName string) (*postgresqlDS, error) {
	// fast path
	r.RLock()
	ds, ok := r.pgCache[dsName]
	r.RUnlock()
	if ok {
		return ds, nil
	}

	// slow path
	r.Lock()
	defer r.Unlock()

	// we do another lookup to make sure any other goroutine didn't create it
	if ds, ok := r.pgCache[dsName]; ok {
		return ds, nil
	}

	for name, info := range manifestdata.GetManifest().Connections {
		if name != dsName {
			continue
		}

		if info.ConnectionType() != manifest.ConnectionTypePostgresql {
			return nil, fmt.Errorf("[%s] is not a postgresql connection", dsName)
		}

		conf := info.(manifest.PostgresqlConnectionInfo)
		if conf.ConnStr == "" {
			return nil, fmt.Errorf("postgresql connection [%s] has empty connString", dsName)
		}

		fullConnStr, err := secrets.ApplySecretsToString(ctx, info, conf.ConnStr)
		if err != nil {
			return nil, fmt.Errorf("failed to apply secrets to connection string for connection [%s]: %w", dsName, err)
		}

		dbpool, err := pgxpool.New(ctx, fullConnStr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres connection [%s]: %w", dsName, err)
		}

		r.pgCache[dsName] = &postgresqlDS{pool: dbpool}
		return r.pgCache[dsName], nil
	}

	return nil, fmt.Errorf("postgresql connection [%s] not found", dsName)
}
