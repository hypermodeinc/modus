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
	"github.com/hypermodeinc/modus/runtime/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresqlDS struct {
	pool *pgxpool.Pool
}

func (ds *postgresqlDS) Shutdown() {
	ds.pool.Close()
}

func (ds *postgresqlDS) query(ctx context.Context, stmt string, params []any, execOnly bool) (*dbResponse, error) {

	tx, err := ds.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error setting up a new tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			logger.Warn(ctx).Err(err).Msg("Error rolling back transaction.")
			return
		}
	}()

	// TODO: what if connection times out and we need to retry

	response := new(dbResponse)
	if execOnly {
		if ct, err := tx.Exec(ctx, stmt, params...); err != nil {
			return nil, err
		} else {
			response.RowsAffected = uint32(ct.RowsAffected())
		}
	} else {
		rows, err := tx.Query(ctx, stmt, params...)
		if err != nil {
			return nil, err
		}

		if data, err := pgx.CollectRows(rows, pgx.RowToMap); err != nil {
			return nil, err
		} else {
			response.Result = data
		}

		response.RowsAffected = uint32(rows.CommandTag().RowsAffected())
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return response, nil
}

func newPostgresqlDS(ctx context.Context, dsName string) (*postgresqlDS, error) {
	connStr, err := getConnectionString(ctx, dsName, manifest.ConnectionTypePostgresql)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql database [%s]: %w", dsName, err)
	}

	return &postgresqlDS{pool}, nil
}
