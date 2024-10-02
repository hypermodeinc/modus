/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sqlclient

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresqlDS struct {
	pool *pgxpool.Pool
}

func (ds *postgresqlDS) query(ctx context.Context, stmt string, params []any) (*dbResponse, error) {

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
	rows, err := tx.Query(ctx, stmt, params...)
	if err != nil {
		return nil, err
	}

	data, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		return nil, err
	}

	rowsAffected := uint32(rows.CommandTag().RowsAffected())

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	response := &dbResponse{
		// Error: "",
		Result:       data,
		RowsAffected: rowsAffected,
	}

	return response, nil
}
