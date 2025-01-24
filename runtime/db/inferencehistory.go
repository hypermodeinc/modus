/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/metrics"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modusdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var globalRuntimePostgresWriter *runtimePostgresWriter = &runtimePostgresWriter{
	dbpool: nil,
	buffer: make(chan inferenceHistory, chanSize),
	quit:   make(chan struct{}),
	done:   make(chan struct{}),
}

type Plugin struct {
	Gid        uint64 `json:"gid,omitempty"`
	Id         string `json:"id,omitempty" db:"constraint=unique"`
	Name       string `json:"name,omitempty"`
	Version    string `json:"version,omitempty"`
	Language   string `json:"language,omitempty"`
	SdkVersion string `json:"sdk_version,omitempty"`
	BuildId    string `json:"build_id,omitempty"`
	BuildTime  string `json:"build_time,omitempty"`
	GitRepo    string `json:"git_repo,omitempty"`
	GitCommit  string `json:"git_commit,omitempty"`
}

type Inference struct {
	Gid        uint64 `json:"gid,omitempty"`
	Id         string `json:"id,omitempty" db:"constraint=unique"`
	ModelHash  string `json:"model_hash,omitempty"`
	Input      string `json:"input,omitempty"`
	Output     string `json:"output,omitempty"`
	StartedAt  string `json:"started_at,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
	Function   string `json:"function,omitempty"`
	Plugin     Plugin `json:"plugin,omitempty"`
}

const batchSize = 100
const chanSize = 10000

const pluginsTable = "plugins"
const inferencesTable = "inferences"

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
	var initErr error
	w.once.Do(func() {
		var connStr string
		var err error
		if secrets.HasSecret("MODUS_DB") {
			connStr, err = secrets.GetSecretValue("MODUS_DB")
		} else if secrets.HasSecret("HYPERMODE_METADATA_DB") {
			// fallback to old secret name
			// TODO: remove this after the transition is complete
			connStr, err = secrets.GetSecretValue("HYPERMODE_METADATA_DB")
		} else {
			return
		}

		if err != nil {
			initErr = err
			return
		}

		if pool, err := pgxpool.New(ctx, connStr); err != nil {
			initErr = err
		} else {
			w.dbpool = pool
		}
	})

	if w.dbpool != nil {
		return w.dbpool, nil
	} else if initErr != nil {
		return nil, initErr
	} else {
		return nil, errDbNotConfigured
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

func WritePluginInfo(ctx context.Context, plugin *plugins.Plugin) {

	if app.IsDevEnvironment() {
		err := writePluginInfoToModusdb(plugin)
		if err != nil {
			logDbWarningOrError(ctx, err, "Plugin info not written to ModusDB.")
		}
		return
	}

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

func writePluginInfoToModusdb(plugin *plugins.Plugin) error {
	if GlobalModusDbEngine == nil {
		return nil
	}
	_, _, err := modusdb.Create[Plugin](GlobalModusDbEngine, Plugin{
		Id:         plugin.Id,
		Name:       plugin.Metadata.Name(),
		Version:    plugin.Metadata.Version(),
		Language:   plugin.Language.Name(),
		SdkVersion: plugin.Metadata.SdkVersion(),
		BuildId:    plugin.Metadata.BuildId,
		BuildTime:  plugin.Metadata.BuildTime,
		GitRepo:    plugin.Metadata.GitRepo,
		GitCommit:  plugin.Metadata.GitCommit,
	})
	return err
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

func WriteInferenceHistoryToDB(ctx context.Context, batch []inferenceHistory) {
	if len(batch) == 0 {
		return
	}

	if app.IsDevEnvironment() {
		err := writeInferenceHistoryToModusDb(batch)
		if err != nil {
			logDbWarningOrError(ctx, err, "Inference history not written to ModusDB.")
		}
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

func writeInferenceHistoryToModusDb(batch []inferenceHistory) error {
	if GlobalModusDbEngine == nil {
		return nil
	}
	for _, data := range batch {
		input, output, err := data.getJson()
		if err != nil {
			return err
		}
		var funcStr string
		var pluginId string
		if data.function == nil {
			funcStr = ""
		} else {
			funcStr = *data.function
		}
		if data.pluginId == nil {
			pluginId = ""
		} else {
			pluginId = *data.pluginId
		}
		_, _, err = modusdb.Create[Inference](GlobalModusDbEngine, Inference{
			Id:         utils.GenerateUUIDv7(),
			ModelHash:  data.model.Hash(),
			Input:      string(input),
			Output:     string(output),
			StartedAt:  data.start.Format(time.RFC3339),
			DurationMs: data.end.Sub(data.start).Milliseconds(),
			Function:   funcStr,
			Plugin: Plugin{
				Id: pluginId,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func QueryPlugins() ([]Plugin, error) {
	if GlobalModusDbEngine == nil {
		return nil, nil
	}
	_, plugins, err := modusdb.Query[Plugin](GlobalModusDbEngine, modusdb.QueryParams{})
	return plugins, err
}

func QueryInferences() ([]Inference, error) {
	if GlobalModusDbEngine == nil {
		return nil, nil
	}
	_, inferences, err := modusdb.Query[Inference](GlobalModusDbEngine, modusdb.QueryParams{})
	return inferences, err
}
