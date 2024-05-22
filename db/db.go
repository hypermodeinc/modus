package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"hmruntime/logger"
	"hmruntime/manifest"
	"hmruntime/metrics"
	"hmruntime/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var globalInferenceWriter *inferenceWriter = &inferenceWriter{}

const batchSize = 100
const chanSize = 10000
const inferencesTable = "inferences"

const inferenceRefresherInterval = 5 * time.Second

type inferenceWriter struct {
	dbpool *pgxpool.Pool
	buffer chan inferenceHistory
	quit   chan struct{}
	done   chan struct{}
}

type inferenceHistory struct {
	model  manifest.Model
	input  any
	output any
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

func WriteInferenceHistory(model manifest.Model, input, output any, start, end time.Time) {
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
			WriteInferenceHistoryToDB(ctx, batch[:batchIndex])
			close(w.done)
			return
		}
	}
}

func WriteInferenceHistoryToDB(ctx context.Context, batch []inferenceHistory) {
	if len(batch) == 0 {
		return
	}
	err := WithTx(ctx, func(tx pgx.Tx) error {
		b := &pgx.Batch{}
		for _, data := range batch {
			input, err := utils.JsonSerialize(data.input)
			if err != nil {
				return err
			}
			output, err := utils.JsonSerialize(data.output)
			if err != nil {
				return err
			}
			query := fmt.Sprintf("INSERT INTO %s (id, model_hash, input, output, started_at, duration_ms) VALUES ($1, $2, $3, $4, $5, $6)", inferencesTable)
			args := []any{utils.GeneratUUID(), data.model.Hash(), input, output, data.start, data.end.Sub(data.start).Milliseconds()}
			b.Queue(query, args...)
		}

		br := tx.SendBatch(ctx, b)
		defer br.Close()

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
		logger.Err(ctx, err).Msg("Error writing to inference history database")
	}

}

func Initialize(ctx context.Context) {
	connStr := os.Getenv("HYPERMODE_METADATA_DB")

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
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
