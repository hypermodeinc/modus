package db

import (
	"context"
	"fmt"
	"hmruntime/logger"
	"hmruntime/manifest"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool
var dbWrites = make(chan InferenceHistory, 10000)
var dbBuffer = make([]InferenceHistory, 0)

type InferenceHistory struct {
	Model  manifest.Model
	Input  string
	Output string
	Start  time.Time
	End    time.Time
}

func WriteInferenceHistory(model manifest.Model, input, output string, start, end time.Time) {
	dbBuffer = append(dbBuffer, InferenceHistory{
		Model:  model,
		Input:  input,
		Output: output,
		Start:  start,
		End:    end,
	})
}

func AsyncFlushInferenceHistory() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			for len(dbBuffer) > 0 {
				select {
				case dbWrites <- dbBuffer[0]:
					dbBuffer = dbBuffer[1:]
				default:
					// dbWrites is full, wait before trying again
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
}

func WriteInferenceHistoryToDB() {
	go func() {
		for data := range dbWrites {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			table := os.Getenv("NAMESPACE")
			if table == "" {
				table = "local_instance"
			}
			version := "1"
			err := WithTx(ctx, func(tx pgx.Tx) error {
				query := fmt.Sprintf(
					"INSERT INTO %s (model_name, model_task, source_model, model_provider, model_host, model_version, model_hash, input, output, started_at, ended_at) "+
						"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", table)
				_, err := tx.Exec(ctx, query, data.Model.Name, data.Model.Task, data.Model.SourceModel, data.Model.Provider, data.Model.Host,
					version, data.Model.Hash(), data.Input, data.Output, data.Start, data.End)
				return err
			})

			if err != nil {
				// Handle error
				fmt.Println("Error writing to inference history database:", err)
			}
			cancel()
		}
	}()

}

func Initialize(ctx context.Context) {
	connStr := os.Getenv("HYPERMODE_CLUSTER_DB")

	var err error
	dbpool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Warn(ctx).Err(err).Msg("Database pool initialization failed.")
	}
	AsyncFlushInferenceHistory()
	WriteInferenceHistoryToDB()
}

func GetInferenceHistoryDB() *pgxpool.Pool {
	return dbpool
}

func GetTx(ctx context.Context) (pgx.Tx, error) {
	if dbpool == nil {
		logger.Warn(ctx).Msg("Database pool is not initialized. Inference history will not be saved")
		return nil, nil
	}
	return dbpool.Begin(ctx)
}

func WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}
	if tx == nil {
		return nil
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("Failed to rollback transaction")
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
