/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/metrics"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/hypermodeAI/manifest"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var globalRuntimePostgresWriter *runtimePostgresWriter = &runtimePostgresWriter{
	dbpool: nil,
	buffer: make(chan inferenceHistory, chanSize),
	quit:   make(chan struct{}),
	done:   make(chan struct{}),
}

var errDbNotConfigured = errors.New("database not configured")

const batchSize = 100
const chanSize = 10000

const pluginsTable = "plugins"
const inferencesTable = "inferences"
const collectionTextsTable = "collection_texts"
const collectionVectorsTable = "collection_vectors"

const inferenceRefresherInterval = 5 * time.Second

type runtimePostgresWriter struct {
	dbpool *pgxpool.Pool
	buffer chan inferenceHistory
	quit   chan struct{}
	done   chan struct{}
	mu     sync.RWMutex
}

type inferenceHistory struct {
	model    *manifest.ModelInfo
	input    any
	output   any
	start    time.Time
	end      time.Time
	pluginId *string
	function *string
}

func (w *runtimePostgresWriter) GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	w.mu.RLock()
	if w.dbpool != nil {
		defer w.mu.RUnlock()
		return w.dbpool, nil
	}
	w.mu.RUnlock()

	w.mu.Lock()
	defer w.mu.Unlock()

	var connStr string
	var err error
	if secrets.HasSecret("MODUS_DB") {
		connStr, err = secrets.GetSecretValue("MODUS_DB")
	} else if secrets.HasSecret("HYPERMODE_METADATA_DB") {
		// fallback to old secret name
		// TODO: remove this after the transition is complete
		connStr, err = secrets.GetSecretValue("HYPERMODE_METADATA_DB")
	} else {
		return nil, errDbNotConfigured
	}
	if err != nil {
		return nil, err
	}

	if pool, err := pgxpool.New(ctx, connStr); err != nil {
		return nil, err
	} else {
		w.dbpool = pool
		return pool, nil
	}
}

func (w *runtimePostgresWriter) Write(data inferenceHistory) {
	select {
	case w.buffer <- data:
	default:
		metrics.DroppedInferencesNum.Inc()
	}
}

func (w *runtimePostgresWriter) worker(ctx context.Context) {
	var batchIndex int
	var batch [batchSize]inferenceHistory
	timer := time.NewTimer(inferenceRefresherInterval)
	defer timer.Stop()

	for {
		select {
		case data := <-w.buffer:
			batch[batchIndex] = data
			batchIndex++
			if batchIndex == batchSize {
				WriteInferenceHistoryToDB(ctx, batch[:batchSize])
				batchIndex = 0

				// we need to drain the timer channel to prevent the timer from firing
				if !timer.Stop() {
					<-timer.C
				}
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

func Stop(ctx context.Context) {
	pool, _ := globalRuntimePostgresWriter.GetPool(ctx)
	if pool == nil {
		return
	}

	close(globalRuntimePostgresWriter.quit)
	<-globalRuntimePostgresWriter.done
	pool.Close()
}

func WritePluginInfo(ctx context.Context, plugin *plugins.Plugin) {

	err := WithTx(ctx, func(tx pgx.Tx) error {

		// Check if the plugin is already in the database
		// If so, update the ID to match
		id, err := getPluginId(ctx, tx, plugin.Metadata.BuildId)
		if err != nil {
			return err
		}
		if id != "" {
			plugin.Id = id
			return nil
		}

		// Insert the plugin info - still check for conflicts, in case another instance of the service is running
		query := fmt.Sprintf(`INSERT INTO %s
(id, name, version, language, sdk_version, build_id, build_time, git_repo, git_commit)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (build_id) DO NOTHING`,
			pluginsTable)

		ct, err := tx.Exec(ctx, query,
			plugin.Id,
			plugin.Metadata.Name(),
			utils.NilIfEmpty(plugin.Metadata.Version()),
			plugin.Language.Name(),
			plugin.Metadata.SdkVersion(),
			plugin.Metadata.BuildId,
			plugin.Metadata.BuildTime,
			utils.NilIfEmpty(plugin.Metadata.GitRepo),
			utils.NilIfEmpty(plugin.Metadata.GitCommit),
		)
		if err != nil {
			return err
		}

		if ct.RowsAffected() == 0 {
			// Edge case - the plugin is now in the database, but we didn't insert it
			// It must have been inserted by another instance of the service
			// Get the ID set by the other instance
			id, err := getPluginId(ctx, tx, plugin.Metadata.BuildId)
			if err != nil {
				return err
			}
			plugin.Id = id
		}

		return nil
	})

	if err != nil {
		logDbWarningOrError(ctx, err, "Plugin info not written to database.")
	}
}

func logDbWarningOrError(ctx context.Context, err error, msg string) {
	if _, ok := err.(*pgconn.ConnectError); ok {
		logger.Warn(ctx).Err(err).Msgf("Database connection error. %s", msg)
	} else if errors.Is(err, errDbNotConfigured) {
		logger.Warn(ctx).Msgf("Database has not been configured. %s", msg)
	} else {
		logger.Err(ctx, err).Msg(msg)
	}
}

func getPluginId(ctx context.Context, tx pgx.Tx, buildId string) (string, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE build_id = $1", pluginsTable)
	rows, err := tx.Query(ctx, query, buildId)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Next() {
		var id string
		err := rows.Scan(&id)
		return id, err
	}
	return "", nil
}

func WriteInferenceHistory(ctx context.Context, model *manifest.ModelInfo, input, output any, start, end time.Time) {
	var pluginId *string
	if plugin, ok := plugins.GetPluginFromContext(ctx); ok {
		pluginId = &plugin.Id
	}

	var function *string
	if functionName, ok := ctx.Value(utils.FunctionNameContextKey).(string); ok {
		function = &functionName
	}

	globalRuntimePostgresWriter.Write(inferenceHistory{
		model:    model,
		input:    input,
		output:   output,
		start:    start,
		end:      end,
		pluginId: pluginId,
		function: function,
	})
}

func GetUniqueNamespaces(ctx context.Context, collectionName string) ([]string, error) {
	var namespaces []string
	err := WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT DISTINCT namespace FROM %s WHERE collection = $1", collectionTextsTable)
		rows, err := tx.Query(ctx, query, collectionName)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var namespace string
			if err := rows.Scan(&namespace); err != nil {
				return err
			}
			namespaces = append(namespaces, namespace)
		}

		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return namespaces, nil
}

func WriteCollectionTexts(ctx context.Context, collectionName, namespace string, keys, texts []string, labelsArr [][]string) ([]int64, error) {
	if len(labelsArr) != 0 && len(keys) != len(labelsArr) {
		return nil, errors.New("if labels is not empty, it must have the same length as keys")
	}

	if len(keys) != len(texts) {
		return nil, errors.New("keys and texts must have the same length")
	}

	ids := make([]int64, len(keys))
	err := WithTx(ctx, func(tx pgx.Tx) error {
		// Delete any existing rows that match the collectionName and keys
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE collection = $1 AND namespace = $2 AND key = ANY($3)", collectionTextsTable)
		_, err := tx.Exec(ctx, deleteQuery, collectionName, namespace, keys)
		if err != nil {
			return err
		}

		// Insert the new rows
		if len(labelsArr) == 0 {
			query := fmt.Sprintf("INSERT INTO %s (collection, namespace, key, text) VALUES ($1, $2, unnest($3::text[]), unnest($4::text[])) RETURNING id", collectionTextsTable)
			rows, err := tx.Query(ctx, query, collectionName, namespace, keys, texts)
			if err != nil {
				return err
			}
			defer rows.Close()

			i := 0
			for rows.Next() {
				if err := rows.Scan(&ids[i]); err != nil {
					return err
				}
				i++
			}
			if err := rows.Err(); err != nil {
				return err
			}
		} else {
			for i := range keys {
				query := fmt.Sprintf("INSERT INTO %s (collection, namespace, key, text, labels) VALUES ($1, $2, $3, $4, $5) RETURNING id", collectionTextsTable)
				err := tx.QueryRow(ctx, query, collectionName, namespace, keys[i], texts[i], labelsArr[i]).Scan(&ids[i])
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return ids, nil
}

func WriteCollectionText(ctx context.Context, collectionName, namespace, key, text string, labels []string) (id int64, err error) {
	err = WithTx(ctx, func(tx pgx.Tx) error {
		// Delete any existing rows that match the collectionName and key
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE collection = $1 AND namespace = $2 AND key = $3", collectionTextsTable)
		_, err := tx.Exec(ctx, deleteQuery, collectionName, namespace, key)
		if err != nil {
			return err
		}

		// Insert the new row
		if len(labels) == 0 {
			labels = nil
		}
		query := fmt.Sprintf("INSERT INTO %s (collection, namespace, key, text, labels) VALUES ($1, $2, $3, $4, unnest($5::text[])) RETURNING id", collectionTextsTable)
		row := tx.QueryRow(ctx, query, collectionName, namespace, key, text, labels)
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

func DeleteCollectionText(ctx context.Context, collectionName, namespace, key string) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("DELETE FROM %s WHERE collection = $1 AND namespace = $2 AND key = $3", collectionTextsTable)
		_, err := tx.Exec(ctx, query, collectionName, namespace, key)
		if err != nil {
			return err
		}
		return nil
	})
}

func WriteCollectionVectors(ctx context.Context, searchMethodName string, textIds []int64, vectors [][]float32) ([]int64, []string, error) {
	if len(textIds) != len(vectors) {
		return nil, nil, errors.New("textIds and vectors must have the same length")
	}

	vectorIds := make([]int64, len(textIds))
	keys := make([]string, len(textIds))
	err := WithTx(ctx, func(tx pgx.Tx) error {
		// Delete any existing rows that match the searchMethodName and textIds
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE search_method = $1 AND text_id = ANY($2)", collectionVectorsTable)
		_, err := tx.Exec(ctx, deleteQuery, searchMethodName, textIds)
		if err != nil {
			return err
		}

		// Insert the new rows
		query := fmt.Sprintf("INSERT INTO %s (search_method, text_id, vector) VALUES ($1, $2, $3::real[]) RETURNING id", collectionVectorsTable)
		for i, textId := range textIds {
			vector := vectors[i]
			row := tx.QueryRow(ctx, query, searchMethodName, textId, vector)
			var id int64
			if err := row.Scan(&id); err != nil {
				return err
			}
			vectorIds[i] = id
		}

		query = "SELECT key FROM collection_texts WHERE id = ANY($1)"
		rows, err := tx.Query(ctx, query, textIds)
		if err != nil {
			return err
		}
		defer rows.Close()

		i := 0
		for rows.Next() {
			if err := rows.Scan(&keys[i]); err != nil {
				return err
			}
			i++
		}
		if err := rows.Err(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return vectorIds, keys, nil
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

func DeleteCollectionVectors(ctx context.Context, collectionName, searchMethodName, namespace string) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf(`
		DELETE FROM %s cv 
		USING %s ct 
		WHERE ct.id = cv.text_id 
		AND ct.collection = $1 
		AND cv.search_method = $2
		AND ct.namespace = $3`,
			collectionVectorsTable, collectionTextsTable)
		_, err := tx.Exec(ctx, query, collectionName, searchMethodName, namespace)
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

func QueryCollectionTextsFromCheckpoint(ctx context.Context, collection, namespace string, textCheckpointId int64) ([]int64, []string, []string, [][]string, error) {
	var textIds []int64
	var keys []string
	var texts []string
	var labelsArr [][]string
	err := WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT id, key, text, labels FROM %s WHERE id > $1 AND collection = $2 AND namespace = $3", collectionTextsTable)
		rows, err := tx.Query(ctx, query, textCheckpointId, collection, namespace)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var id int64
			var key string
			var text string
			var labels []string
			if err := rows.Scan(&id, &key, &text, &labels); err != nil {
				return err
			}
			textIds = append(textIds, id)
			keys = append(keys, key)
			texts = append(texts, text)
			labelsArr = append(labelsArr, labels)

		}

		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, nil, nil, nil, err
	}

	return textIds, keys, texts, labelsArr, nil
}

func QueryCollectionVectorsFromCheckpoint(ctx context.Context, collectionName, searchMethodName, namespace string, vecCheckpointId int64) ([]int64, []int64, []string, [][]float32, error) {
	var textIds []int64
	var vectorIds []int64
	var keys []string
	var vectors [][]float32
	err := WithTx(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf(`SELECT ct.id, cv.id, ct.key, cv.vector 
                  FROM %s cv 
                  JOIN %s ct ON cv.text_id = ct.id 
                  WHERE cv.id > $1 AND ct.collection = $2 AND cv.search_method = $3 AND ct.namespace = $4`, collectionVectorsTable, collectionTextsTable)
		rows, err := tx.Query(ctx, query, vecCheckpointId, collectionName, searchMethodName, namespace)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var textId int64
			var vectorId int64
			var key string
			var vector []float32
			if err := rows.Scan(&textId, &vectorId, &key, &vector); err != nil {
				return err
			}
			textIds = append(textIds, textId)
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
		return nil, nil, nil, nil, err
	}

	return textIds, vectorIds, keys, vectors, nil
}

func WriteInferenceHistoryToDB(ctx context.Context, batch []inferenceHistory) {
	if len(batch) == 0 {
		return
	}

	err := WithTx(ctx, func(tx pgx.Tx) error {
		b := &pgx.Batch{}
		for _, data := range batch {
			input, output, err := data.getJson()
			if err != nil {
				return err
			}
			query := fmt.Sprintf(`INSERT INTO %s
(id, model_hash, input, output, started_at, duration_ms, plugin_id, function)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`, inferencesTable)
			args := []any{
				utils.GenerateUUIDv7(),
				data.model.Hash(),
				input,
				output,
				data.start,
				data.end.Sub(data.start).Milliseconds(),
				data.pluginId,
				data.function,
			}
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
		logDbWarningOrError(ctx, err, "Inference history not written to database.")
	}
}

func Initialize(ctx context.Context) {
	// this will initialize the pool and start the worker
	_, err := globalRuntimePostgresWriter.GetPool(ctx)
	if err != nil {
		logger.Warn(ctx).Err(err).Msg("Metadata database is not available.")
	}
	go globalRuntimePostgresWriter.worker(ctx)
}

func GetTx(ctx context.Context) (pgx.Tx, error) {
	pool, err := globalRuntimePostgresWriter.GetPool(ctx)
	if err != nil {
		return nil, err
	}

	return pool.Begin(ctx)
}

func WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	span, ctx := utils.NewSentrySpanForCallingFunc(ctx)
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
