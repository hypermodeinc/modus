package db

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"hmruntime/logger"
	"hmruntime/metrics"
	"hmruntime/plugins"
	"hmruntime/secrets"
	"hmruntime/utils"

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
	connStr, err := secrets.GetSecretValue(ctx, "HYPERMODE_METADATA_DB")
	if connStr == "" {
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errDbNotConfigured, err)
		} else {
			return nil, errDbNotConfigured
		}
	} else if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	w.dbpool = pool
	return pool, nil
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
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
		const query = `INSERT INTO plugins
(id, name, version, language, sdk_version, build_id, build_time, git_repo, git_commit)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (build_id) DO NOTHING`

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
	const query = "SELECT id FROM plugins WHERE build_id = $1"
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
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	var pluginId *string
	if plugin, ok := ctx.Value(utils.PluginContextKey).(*plugins.Plugin); ok {
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

func getCollectionId(ctx context.Context, tx pgx.Tx, collectionName string) (int32, error) {
	var id int32
	const query = "SELECT id FROM collections WHERE c.name = $1"
	err := tx.QueryRow(ctx, query, collectionName).Scan(&id)
	if err == pgx.ErrNoRows {
		const insertQuery = "INSERT INTO collections (name) VALUES ($1) RETURNING id"
		if err := tx.QueryRow(ctx, insertQuery, collectionName).Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, err
}

func getNamespaceId(ctx context.Context, tx pgx.Tx, collection, namespace string) (int32, error) {
	const query = `
		SELECT n.id
		FROM collection_namespaces n
		JOIN collections c ON c.id = n.collection_id
		WHERE c.name = $1 AND n.name = $2
	`
	var id int32
	err := tx.QueryRow(ctx, query, collection, namespace).Scan(&id)
	if err == pgx.ErrNoRows {
		cId, err := getCollectionId(ctx, tx, collection)
		if err != nil {
			return 0, err
		}
		const insertQuery = "INSERT INTO collection_namespaces (collection_id, name) VALUES ($1, $2) RETURNING id"
		if err := tx.QueryRow(ctx, insertQuery, cId, namespace).Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, err
}

func GetUniqueNamespaces(ctx context.Context, collectionName string) ([]string, error) {
	var namespaces []string
	err := WithTx(ctx, func(tx pgx.Tx) error {
		const query = `
			SELECT n.name
			FROM collection_namespaces n
			JOIN collections c on c.id = n.collection_id
			WHERE c.name = $1
		`
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

func WriteCollectionTexts(ctx context.Context, collectionName, namespace string, keys, texts []string, labels [][]string) ([]int64, error) {
	if len(labels) != 0 && len(keys) != len(labels) {
		return nil, errors.New("if labels is not empty, it must have the same length as keys")
	}

	if len(keys) != len(texts) {
		return nil, errors.New("keys and texts must have the same length")
	}

	ids := make([]int64, len(keys))
	err := WithTx(ctx, func(tx pgx.Tx) error {

		// Get the namespace ID
		nsId, err := getNamespaceId(ctx, tx, collectionName, namespace)
		if err != nil {
			return err
		}

		// Delete any existing rows for the given keys
		const deleteQuery = "DELETE FROM collection_texts WHERE namespace_id = $1 AND key = ANY($2)"
		if _, err := tx.Exec(ctx, deleteQuery, nsId, keys); err != nil {
			return err
		}

		// Insert the new rows
		var rows pgx.Rows
		if len(labels) == 0 {
			const query = "INSERT INTO collection_texts (namespace_id, key, text) VALUES ($1, unnest($2::text[]), unnest($3::text[])) RETURNING id"
			rows, err = tx.Query(ctx, query, nsId, keys, texts)
		} else {
			const query = "INSERT INTO collection_texts (namespace_id, key, text, labels) VALUES ($1, unnest($2::text[]), unnest($3::text[]), unnest_nd_1d($4::text[][])) RETURNING id"
			rows, err = tx.Query(ctx, query, nsId, keys, texts, labels)
		}
		if err != nil {
			return err
		}
		defer rows.Close()

		for i := 0; rows.Next(); i++ {
			if err := rows.Scan(&ids[i]); err != nil {
				return err
			}
		}
		if err := rows.Err(); err != nil {
			return err
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

		// Get the namespace ID
		nsId, err := getNamespaceId(ctx, tx, collectionName, namespace)
		if err != nil {
			return err
		}

		// Delete any existing rows for the given key
		const deleteQuery = "DELETE FROM collection_texts WHERE namespace_id = $1 AND key = $2"
		if _, err := tx.Exec(ctx, deleteQuery, nsId, key); err != nil {
			return err
		}

		// Insert the new row
		if len(labels) == 0 {
			labels = nil
		}
		const query = "INSERT INTO collection_texts (namespace_id, key, text, labels) VALUES ($1, $2, $3, unnest($4::text[])) RETURNING id"
		row := tx.QueryRow(ctx, query, nsId, key, text, labels)
		return row.Scan(&id)
	})

	if err != nil {
		return 0, err
	}
	return id, nil
}

func DeleteCollection(ctx context.Context, collectionName string) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		// Deletes the collection, including all associated namespaces, texts and vectors (via cascading delete)
		const query = `
			DELETE FROM collection_namespaces n
			USING collections c
			WHERE c.name = $1 AND n.name = $2
				AND c.id = n.collection_id
		`
		if _, err := tx.Exec(ctx, query, collectionName); err != nil {
			return err
		}
		return nil
	})
}

func DeleteNamespace(ctx context.Context, collectionName, namespace string) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		// Deletes the namespace, including all associated texts and vectors (via cascading delete)
		const query = "DELETE FROM collections WHERE c.name = $1"
		if _, err := tx.Exec(ctx, query, collectionName, namespace); err != nil {
			return err
		}
		return nil
	})
}

func DeleteCollectionText(ctx context.Context, collectionName, namespace, key string) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		const query = `
			DELETE FROM collection_texts t
			USING collections c, collection_namespaces n
			WHERE c.name = $1 AND n.name = $2 AND t.key = $3
				AND c.id = n.collection_id
				AND n.id = t.namespace_id
		`
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
		const deleteQuery = "DELETE FROM collection_vectors WHERE search_method = $1 AND text_id = ANY($2)"
		if _, err := tx.Exec(ctx, deleteQuery, searchMethodName, textIds); err != nil {
			return err
		}

		// Insert the new rows
		const insertQuery = "INSERT INTO collection_vectors (search_method, text_id, vector) VALUES ($1, $2, $3::real[]) RETURNING id"
		for i, textId := range textIds {
			vector := vectors[i]
			row := tx.QueryRow(ctx, insertQuery, searchMethodName, textId, vector)
			var id int64
			if err := row.Scan(&id); err != nil {
				return err
			}
			vectorIds[i] = id
		}

		const query = "SELECT key FROM collection_texts WHERE id = ANY($1)"
		rows, err := tx.Query(ctx, query, textIds)
		if err != nil {
			return err
		}
		defer rows.Close()

		for i := 0; rows.Next(); i++ {
			if err := rows.Scan(&keys[i]); err != nil {
				return err
			}
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
		const deleteQuery = "DELETE FROM collection_vectors WHERE search_method = $1 AND text_id = $2"
		if _, err := tx.Exec(ctx, deleteQuery, searchMethodName, textId); err != nil {
			return err
		}

		// Insert the new row
		const insertQuery = "INSERT INTO collection_vectors (search_method, text_id, vector) VALUES ($1, $2, $3) RETURNING id"
		row := tx.QueryRow(ctx, insertQuery, searchMethodName, textId, vector)
		if err := row.Scan(&vectorId); err != nil {
			return err
		}

		const query = "SELECT key FROM collection_texts WHERE id = $1"
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
		const query = `
			DELETE FROM collection_vectors v
			USING collection_texts t, collection_namespaces n, collections c
			WHERE t.id = v.text_id
				AND n.id = t.namespace_id
				AND c.id = n.collection_id
				AND c.name = $1
				AND n.name = $2
				AND v.search_method = $3
		`
		if _, err := tx.Exec(ctx, query, collectionName, namespace, searchMethodName); err != nil {
			return err
		}
		return nil
	})
}

func DeleteCollectionVector(ctx context.Context, searchMethodName string, textId int64) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
		const query = "DELETE FROM collection_vectors WHERE search_method = $1 AND text_id = $2"
		if _, err := tx.Exec(ctx, query, searchMethodName, textId); err != nil {
			return err
		}
		return nil
	})
}

func QueryCollectionTextsFromCheckpoint(ctx context.Context, collection, namespace string, textCheckpointId int64) ([]int64, []string, []string, [][]string, error) {
	var textIds []int64
	var keys, texts []string
	var labelsArr [][]string
	err := WithTx(ctx, func(tx pgx.Tx) error {

		nsId, err := getNamespaceId(ctx, tx, collection, namespace)
		if err != nil {
			return err
		}

		const query = "SELECT id, key, text, labels FROM collection_texts WHERE namespace_id = $1 AND id > $2"
		rows, err := tx.Query(ctx, query, nsId, textCheckpointId)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var id int64
			var key, text string
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
	var textIds, vectorIds []int64
	var keys []string
	var vectors [][]float32
	err := WithTx(ctx, func(tx pgx.Tx) error {

		nsId, err := getNamespaceId(ctx, tx, collectionName, namespace)
		if err != nil {
			return err
		}

		const query = `
			SELECT t.id, v.id, t.key, v.vector 
			FROM collection_vectors v 
			JOIN collection_texts t ON t.id = v.text_id
			WHERE t.namespace_id = $1 AND v.search_method = $2 AND v.id > $3
		`
		rows, err := tx.Query(ctx, query, nsId, searchMethodName, vecCheckpointId)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var textId, vectorId int64
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
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	const query = `
		INSERT INTO inferences
		(id, model_hash, input, output, started_at, duration_ms, plugin_id, function)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	err := WithTx(ctx, func(tx pgx.Tx) error {
		b := &pgx.Batch{}
		for _, data := range batch {
			input, output, err := data.getJson()
			if err != nil {
				return err
			}

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
	if _, err := globalRuntimePostgresWriter.GetPool(ctx); err != nil {
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
	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			return
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
