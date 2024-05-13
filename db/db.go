package db

import (
	"context"
	"fmt"
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
	if dbpool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	return dbpool.Begin(ctx)
}

func WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = fn(tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
