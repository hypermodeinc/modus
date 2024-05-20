package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"hmruntime/logger"
	"hmruntime/manifest"
	"hmruntime/metrics"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var globalInferenceWriter *inferenceWriter

const batchSize = 100
const chanSize = 10000

const inferenceRefresherInterval = 5 * time.Second

type inferenceWriter struct {
	dbpool *pgxpool.Pool
	buffer chan inferenceHistory
	quit   chan struct{}
	done   chan struct{}
}

type inferenceHistory struct {
	model  manifest.Model
	input  string
	output string
	start  time.Time
	end    time.Time
}

func (w *inferenceWriter) Write(data inferenceHistory) {
	select {
	case w.buffer <- data:
	default:
		metrics.DroppedInferencesNum.Inc()
	}
}

func Stop() {
	close(globalInferenceWriter.quit)
	<-globalInferenceWriter.done
	globalInferenceWriter.dbpool.Close()
}

func WriteInferenceHistory(model manifest.Model, input, output string, start, end time.Time) {
	globalInferenceWriter.Write(inferenceHistory{
		model:  model,
		input:  input,
		output: output,
		start:  start,
		end:    end,
	})
}

func (w *inferenceWriter) worker(ctx context.Context) {
	var batchIndex int
	var batch [batchSize]inferenceHistory
	timer := time.NewTimer(inferenceRefresherInterval)
	for {
		select {
		case data := <-w.buffer:
			batch[batchIndex] = data
			batchIndex++
			if batchIndex == batchSize {
				WriteInferenceHistoryToDB(ctx, batch[:batchSize])
				batchIndex = 0
				timer.Reset(inferenceRefresherInterval)
			}
		case <-timer.C:
			WriteInferenceHistoryToDB(ctx, batch[:batchIndex])
			batchIndex = 0
			timer.Reset(inferenceRefresherInterval)
		case <-w.quit:
			WriteInferenceHistoryToDB(ctx, batch[:batchSize])
			close(w.done)
			return
		}
	}
}

func WriteInferenceHistoryToDB(ctx context.Context, batch []inferenceHistory) {
	if len(batch) == 0 {
		return
	}
	table := os.Getenv("NAMESPACE")
	if table == "" {
		table = "local_instance"
	}
	err := WithTx(ctx, func(tx pgx.Tx) error {
		if tx == nil {
			return nil
		}
		b := &pgx.Batch{}
		for _, data := range batch {
			query := fmt.Sprintf("INSERT INTO %s (model_hash, input, output, started_at, duration_ms) VALUES ($1, $2, $3, $4, $5)", table)
			args := []interface{}{data.model.Hash(), data.input, data.output, data.start, data.end.Sub(data.start).Milliseconds()}
			b.Queue(query, args...)
		}

		br := tx.SendBatch(ctx, b)
		defer func() {
			if err := br.Close(); err != nil {
				logger.Error(ctx).Err(err).Msg("Error closing batch")
			}
		}()

		for range batch {
			_, err := br.Exec()
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		// Handle error
		logger.Error(ctx).Err(err).Msg("Error writing to inference history database")
	}

}

func Initialize(ctx context.Context) {
	connStr := os.Getenv("HYPERMODE_CLUSTER_DB")

	var err error
	tempDBPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Warn(ctx).Err(err).Msg("Database pool initialization failed.")
	}
	globalInferenceWriter = &inferenceWriter{
		dbpool: tempDBPool,
		buffer: make(chan inferenceHistory, chanSize),
		quit:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	go globalInferenceWriter.worker(ctx)
}

func GetTx(ctx context.Context) (pgx.Tx, error) {
	if globalInferenceWriter.dbpool == nil {
		logger.Warn(ctx).Msg("Database pool is not initialized. Inference history will not be saved")
		return nil, nil
	}
	return globalInferenceWriter.dbpool.Begin(ctx)
}

func WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := GetTx(ctx)
	if err != nil {
		return err
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
