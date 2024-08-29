/*
 * Copyright 2024 Hypermode, Inc.
 */

package sqlclient

import (
	"context"
	"fmt"
	"hypruntime/logger"

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
