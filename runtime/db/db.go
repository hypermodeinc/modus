/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"context"
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errDbNotConfigured = errors.New("database not configured")

type runtimePostgresWriter struct {
	dbpool *pgxpool.Pool
	once   sync.Once
}

func Stop(ctx context.Context) {
	pool, _ := globalRuntimePostgresWriter.getPool(ctx)
	if pool == nil {
		return
	}
	pool.Close()
}

func logDbError(ctx context.Context, err error, msg string) {
	sentryutils.CaptureError(ctx, err, msg)
	if _, ok := err.(*pgconn.ConnectError); ok {
		logger.Error(ctx, err).Msgf("Database connection error. %s", msg)
	} else if errors.Is(err, errDbNotConfigured) {
		if !useModusDB() {
			logger.Error(ctx).Msgf("Database has not been configured. %s", msg)
		}
	} else {
		logger.Error(ctx, err).Bool("user_visible", true).Msg(msg)
	}
}

func Initialize(ctx context.Context) {
	if useModusDB() {
		return
	}

	_, err := globalRuntimePostgresWriter.getPool(ctx)
	if err != nil {
		logger.Warn(ctx, err).Msg("Metadata database is not available.")
	}
}

func GetTx(ctx context.Context) (pgx.Tx, error) {
	pool, err := globalRuntimePostgresWriter.getPool(ctx)
	if err != nil {
		return nil, err
	}

	return pool.Begin(ctx)
}

func WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	span, ctx := sentryutils.NewSpanForCallingFunc(ctx)
	defer span.Finish()

	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			return
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

var _useModusDBOnce sync.Once
var _useModusDB bool

func useModusDB() bool {
	_useModusDBOnce.Do(func() {
		s := os.Getenv("MODUS_USE_MODUSDB")
		if s != "" {
			if value, err := strconv.ParseBool(s); err == nil {
				_useModusDB = value
				return
			}
		}

		_useModusDB = app.IsDevEnvironment()
	})
	return _useModusDB
}
