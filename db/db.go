package db

import (
	"context"
	"fmt"
	"time"

	"hmruntime/config"
	"hmruntime/logger"
	"hmruntime/metrics"
	"hmruntime/secrets"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var globalRuntimePostgresWriter *runtimePostgresWriter = &runtimePostgresWriter{
	dbpool: nil,
	buffer: make(chan inferenceHistory, chanSize),
	quit:   make(chan struct{}),
	done:   make(chan struct{}),
}

const batchSize = 100
const chanSize = 10000
const inferencesTable = "inferences"
const collectionTextsTable = "collection_texts"
const collectionVectorsTable = "collection_vectors"

const inferenceRefresherInterval = 5 * time.Second

type runtimePostgresWriter struct {
	dbpool *pgxpool.Pool
	buffer chan inferenceHistory
	quit   chan struct{}
	done   chan struct{}
}

type inferenceHistory struct {
	model  manifest.ModelInfo
	input  any
	output any
	start  time.Time
	end    time.Time
}

func (h *inferenceHistory) getJson() (input []byte, output []byte, err error) {
	input, err = getInferenceDataJson(h.input)
	if err != nil {
		return nil, nil, err
	}
	output, err = getInferenceDataJson(h.output)
	if err != nil {
		return nil, nil, err
	}
	return input, output, nil
}

func getInferenceDataJson(val any) ([]byte, error) {

	// If the value is a byte slice or string, it must already have been serialized as JSON.
	// It might be formatted, but we don't care because we store in a JSONB column in Postgres,
	// which doesn't preserve formatting.
	switch t := val.(type) {
	case []byte:
		return t, nil
	case string:
		return []byte(t), nil
	}

	// For all other types, we serialize to JSON ourselves.
	bytes, err := utils.JsonSerialize(val)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (w *runtimePostgresWriter) Write(data inferenceHistory) {
	select {
	case w.buffer <- data:
	default:
		metrics.DroppedInferencesNum.Inc()
	}
}

func Stop() {
	close(globalRuntimePostgresWriter.quit)
	<-globalRuntimePostgresWriter.done
	globalRuntimePostgresWriter.dbpool.Close()
}

func WriteInferenceHistory(ctx context.Context, model manifest.ModelInfo, input, output any, start, end time.Time) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()
	globalRuntimePostgresWriter.Write(inferenceHistory{
		model:  model,
		input:  input,
		output: output,
		start:  start,
		end:    end,
	})
}

func WriteCollectionText(ctx context.Context, collectionName, key, text string) (id int64, err error) {
	err = WithTx(ctx, func(tx pgx.Tx) error {
		// Delete any existing rows that match the collectionName and key
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE collection = $1 AND key = $2", collectionTextsTable)
		_, err := tx.Exec(ctx, deleteQuery, collectionName, key)
		if err != nil {
			return err
		}

		// Insert the new row
		query := fmt.Sprintf("INSERT INTO %s (collection, key, text) VALUES ($1, $2, $3) RETURNING id", collectionTextsTable)
		row := tx.QueryRow(ctx, query, collectionName, key, text)
		return row.Scan(&id)
	})

	if err != nil {
		return 0, err
	}
	return id, nil
}

func DeleteCollectionTexts(ctx context.Context, collectionName string) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("DELETE FROM %s WHERE collection = $1", collectionTextsTable)
		_, err := tx.Exec(ctx, query, collectionName)
		if err != nil {
			return err
		}
		return nil
	})
}

func DeleteCollectionText(ctx context.Context, collectionName, key string) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("DELETE FROM %s WHERE collection = $1 AND key = $2", collectionTextsTable)
		_, err := tx.Exec(ctx, query, collectionName, key)
		if err != nil {
			return err
		}
		return nil
	})
}

func WriteCollectionVector(ctx context.Context, searchMethodName string, textId int64, vector []float32) (vectorId int64, key string, err error) {
	err = WithTx(ctx, func(tx pgx.Tx) error {
		// Delete any existing rows that match the searchMethodName and textId
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE search_method = $1 AND text_id = $2", collectionVectorsTable)
		_, err := tx.Exec(ctx, deleteQuery, searchMethodName, textId)
		if err != nil {
			return err
		}

		// Insert the new row
		query := fmt.Sprintf("INSERT INTO %s (search_method, text_id, vector) VALUES ($1, $2, $3) RETURNING id", collectionVectorsTable)
		row := tx.QueryRow(ctx, query, searchMethodName, textId, vector)
		err = row.Scan(&vectorId)
		if err != nil {
			return err
		}

		query = "SELECT key FROM collection_texts WHERE id = $1"
		row = tx.QueryRow(ctx, query, textId)
		return row.Scan(&key)

	})

	if err != nil {
		return 0, "", err
	}
	return vectorId, key, nil
}

func DeleteCollectionVectors(ctx context.Context, collectionName, searchMethodName string) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf(`
		DELETE FROM %s cv 
		USING %s ct 
		WHERE ct.id = cv.text_id 
		AND ct.collection = $1 
		AND cv.search_method = $2`,
			collectionVectorsTable, collectionTextsTable)
		_, err := tx.Exec(ctx, query, collectionName, searchMethodName)
		if err != nil {
			return err
		}
		return nil
	})
}

func DeleteCollectionVector(ctx context.Context, searchMethodName string, textId int64) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("DELETE FROM %s WHERE search_method = $1 AND text_id = $2", collectionVectorsTable)
		_, err := tx.Exec(ctx, query, searchMethodName, textId)
		if err != nil {
			return err
		}
		return nil
	})
}

func QueryCollectionTextsFromCheckpoint(ctx context.Context, collection string, textCheckpointId int64) ([]int64, []string, []string, error) {
	var textIds []int64
	var keys []string
	var texts []string
	err := WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT id, key, text FROM %s WHERE id > $1 AND collection = $2", collectionTextsTable)
		rows, err := tx.Query(ctx, query, textCheckpointId, collection)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var id int64
			var key string
			var text string
			if err := rows.Scan(&id, &key, &text); err != nil {
				return err
			}
			textIds = append(textIds, id)
			keys = append(keys, key)
			texts = append(texts, text)
		}

		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return textIds, keys, texts, nil
}

func QueryCollectionVectorsFromCheckpoint(ctx context.Context, collectionName, searchMethodName string, vecCheckpointId int64) ([]int64, []string, [][]float32, error) {
	var vectorIds []int64
	var keys []string
	var vectors [][]float32
	err := WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf(`SELECT cv.id, ct.key, cv.vector 
                  FROM %s cv 
                  JOIN %s ct ON cv.text_id = ct.id 
                  WHERE cv.id > $1 AND ct.collection = $2 AND cv.search_method = $3`, collectionVectorsTable, collectionTextsTable)
		rows, err := tx.Query(ctx, query, vecCheckpointId, collectionName, searchMethodName)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var vectorId int64
			var key string
			var vector []float32
			if err := rows.Scan(&vectorId, &key, &vector); err != nil {
				return err
			}
			vectorIds = append(vectorIds, vectorId)
			keys = append(keys, key)
			vectors = append(vectors, vector)
		}

		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return vectorIds, keys, vectors, nil
}

func (w *runtimePostgresWriter) worker(ctx context.Context) {
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
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()
	err := WithTx(ctx, func(tx pgx.Tx) error {
		b := &pgx.Batch{}
		for _, data := range batch {
			input, output, err := data.getJson()
			if err != nil {
				return err
			}
			query := fmt.Sprintf("INSERT INTO %s (id, model_hash, input, output, started_at, duration_ms) VALUES ($1, $2, $3, $4, $5, $6)", inferencesTable)
			args := []any{utils.GenerateUUIDV7(), data.model.Hash(), input, output, data.start, data.end.Sub(data.start).Milliseconds()}
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
		if !config.IsDevEnvironment() {
			logger.Err(ctx, err).Msg("Error writing to inference history database")
		}
	}

}

func Initialize(ctx context.Context) {
	connStr, err := secrets.GetSecretValue(ctx, "HYPERMODE_METADATA_DB")
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting database connection string")
		return
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Warn(ctx).Err(err).Msg("Database pool initialization failed.")
	}
	globalRuntimePostgresWriter.dbpool = pool
	go globalRuntimePostgresWriter.worker(ctx)
}

func GetTx(ctx context.Context) (pgx.Tx, error) {
	if globalRuntimePostgresWriter.dbpool == nil {
		logger.Warn(ctx).Msg("Database pool is not initialized. Inference history will not be saved")
		return nil, nil
	}
	return globalRuntimePostgresWriter.dbpool.Begin(ctx)
}

func WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
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
