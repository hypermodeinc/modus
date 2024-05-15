package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"hmruntime/logger"
	"hmruntime/manifest"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var globalInferenceWriter *inferenceWriter

const batchSize = 100
const chanSize = 10000

type inferenceWriter struct {
	dbpool     *pgxpool.Pool
	buffer     chan inferenceHistory
	numDropped int
	quit       chan struct{}
	done       chan struct{}
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
		w.numDropped++
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
	timer := time.NewTimer(5 * time.Second)
	for {
		select {
		case data := <-w.buffer:
			batch[batchIndex] = data
			batchIndex++
			if batchIndex == batchSize {
				w.flush(ctx, batch[:batchSize])
				batchIndex = 0
				timer.Reset(5 * time.Second)
			}
		case <-timer.C:
			w.flush(ctx, batch[:batchIndex])
			batchIndex = 0
			timer.Reset(5 * time.Second)
		case <-w.quit:
			w.flush(ctx, batch[:batchSize])
			close(w.done)
			return
		}
	}
}

func (w *inferenceWriter) flush(ctx context.Context, batch []inferenceHistory) {
	if len(batch) == 0 {
		return
	}
	WriteInferenceHistoryToDB(ctx, batch)
}

func WriteInferenceHistoryToDB(ctx context.Context, batch []inferenceHistory) {
	table := os.Getenv("NAMESPACE")
	if table == "" {
		table = "local_instance"
	}
	err := WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("INSERT INTO %s (model_hash, input, output, started_at, duration_ms) VALUES ", table)
		args := []interface{}{}

		for i, data := range batch {
			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
			if i != len(batch)-1 {
				query += ", "
			}
			args = append(args, data.model.Hash(), data.input, data.output, data.start, data.end.Sub(data.start).Milliseconds())
		}

		_, err := tx.Exec(ctx, query, args...)
		return err
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
