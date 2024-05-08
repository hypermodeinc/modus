package db

import (
	"context"
	"hmruntime/config"
	"hmruntime/logger"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func Initialize() {
	var connStr string
	if config.GetEnvironmentName() == "dev" {
		connStr = "postgresql://postgres:postgres@localhost:5433/my-runtime-db?sslmode=disable"
	} else {
		connStr = os.Getenv("DB_CONN_STR")
	}

	var err error
	dbpool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Fatal(context.Background()).Err(err).Msg("Error connecting to inference history database.")
	}
}

func GetInferenceHistoryDB() *pgxpool.Pool {
	return dbpool
}

func GetTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func WithTx(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}

	err = fn(ctx, tx)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			logger.Error(ctx).Err(rollbackErr).Msg("Error rolling back transaction.")
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
