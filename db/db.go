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
var GlobalInferenceWriter *InferenceWriter

const batchSize = 100
const chanSize = 10000

type InferenceWriter struct {
	buffer     chan InferenceHistory
	numDropped int
	quit       chan struct{}
	done       chan struct{}
}

type InferenceHistory struct {
	Model  manifest.Model
	Input  string
	Output string
	Start  time.Time
	End    time.Time
}

func NewInferenceWriter() {
	GlobalInferenceWriter = &InferenceWriter{
		buffer: make(chan InferenceHistory, chanSize),
		quit:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	go GlobalInferenceWriter.worker()
}

func (w *InferenceWriter) Write(data InferenceHistory) {
	select {
	case w.buffer <- data:
	default:
		w.numDropped++
	}
}

func (w *InferenceWriter) Stop() {
	close(w.quit)
	<-w.done
}

func (w *InferenceWriter) worker() {
	var batchIndex int
	var batch [batchSize]InferenceHistory
	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case data := <-w.buffer:
			batch[batchIndex] = data
			batchIndex++
			if batchIndex == batchSize {
				w.flush(batch[:batchSize])
				batchIndex = 0
			}
		case <-timer.C:
			w.flush(batch[:batchIndex])
			batchIndex = 0
			timer.Reset(10 * time.Second)
		case <-w.quit:
			w.flush(batch[:batchSize])
			close(w.done)
			return
		}
	}
}

func (w *InferenceWriter) flush(batch []InferenceHistory) {
	if len(batch) == 0 {
		return
	}
	WriteInferenceHistoryToDB(batch)
}

func WriteInferenceHistory(model manifest.Model, input, output string, start, end time.Time) {
	GlobalInferenceWriter.Write(InferenceHistory{
		Model:  model,
		Input:  input,
		Output: output,
		Start:  start,
		End:    end,
	})
}

func WriteInferenceHistoryToDB(batch []InferenceHistory) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	table := os.Getenv("NAMESPACE")
	if table == "" {
		table = "local_instance"
	}
	version := "1"
	err := WithTx(ctx, func(tx pgx.Tx) error {
		for _, data := range batch {
			query := fmt.Sprintf(
				"INSERT INTO %s (model_name, model_task, source_model, model_provider, model_host, model_version, model_hash, input, output, started_at, ended_at) "+
					"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", table)
			_, err := tx.Exec(ctx, query, data.Model.Name, data.Model.Task, data.Model.SourceModel, data.Model.Provider, data.Model.Host,
				version, data.Model.Hash(), data.Input, data.Output, data.Start, data.End)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		// Handle error
		fmt.Println("Error writing to inference history database:", err)
	}
	cancel()

}

func Initialize(ctx context.Context) {
	connStr := os.Getenv("HYPERMODE_CLUSTER_DB")

	var err error
	dbpool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Warn(ctx).Err(err).Msg("Database pool initialization failed.")
	}
	NewInferenceWriter()
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
