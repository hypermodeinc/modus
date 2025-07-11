/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/tidwall/gjson"

	"github.com/hypermodeinc/modusgraph"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var globalRuntimePostgresWriter *runtimePostgresWriter = &runtimePostgresWriter{}

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
	Plugin     Plugin `json:"plugin"`
}

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

func (w *runtimePostgresWriter) getPool(ctx context.Context) (*pgxpool.Pool, error) {
	var initErr error
	w.once.Do(func() {
		if !secrets.HasSecret(ctx, "MODUS_DB") {
			return
		}

		connStr, err := secrets.GetSecretValue(ctx, "MODUS_DB")
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
	// which doesn't preserve formatting. For all other types, we serialize to JSON ourselves.

	var result []byte
	switch t := val.(type) {
	case []byte:
		result = t
	case string:
		result = []byte(t)
	default:
		if b, err := utils.JsonSerialize(val); err == nil {
			result = b
		} else {
			return nil, err
		}
	}

	result = utils.SanitizeUTF8(result)

	if !gjson.ValidBytes(result) {
		return nil, fmt.Errorf("invalid JSON data: %s", result)
	}

	return result, nil
}

func WritePluginInfo(ctx context.Context, plugin *plugins.Plugin) {
	var err error
	if useModusDB() {
		err = writePluginInfoToModusDB(ctx, plugin)
	} else {
		err = writePluginInfoToPostgresDB(ctx, plugin)
	}

	if err != nil {
		logDbError(ctx, err, "Plugin info not written to database.")
	}
}

func writePluginInfoToModusDB(ctx context.Context, plugin *plugins.Plugin) error {
	if globalModusDbEngine == nil {
		return nil
	}
	_, _, err := modusgraph.Create(ctx, globalModusDbEngine, Plugin{
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

func writePluginInfoToPostgresDB(ctx context.Context, plugin *plugins.Plugin) error {
	return WithTx(ctx, func(tx pgx.Tx) error {

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

	data := inferenceHistory{
		model:    model,
		input:    input,
		output:   output,
		start:    start,
		end:      end,
		pluginId: pluginId,
		function: function,
	}

	var err error
	if useModusDB() {
		err = writeInferenceHistoryToModusDB(ctx, data)
	} else {
		err = writeInferenceHistoryToPostgresDB(ctx, data)
	}

	if err != nil {
		logDbError(ctx, err, "Inference history not written to database.")
	}
}

func writeInferenceHistoryToModusDB(ctx context.Context, data inferenceHistory) error {
	if globalModusDbEngine == nil {
		return nil
	}

	input, output, err := data.getJson()
	if err != nil {
		return err
	}

	var funcStr string
	var pluginId string
	if data.function != nil {
		funcStr = *data.function
	}
	if data.pluginId != nil {
		pluginId = *data.pluginId
	}

	_, _, err = modusgraph.Create(ctx, globalModusDbEngine, Inference{
		Id:         utils.GenerateUUIDv7(),
		ModelHash:  data.model.Hash(),
		Input:      string(input),
		Output:     string(output),
		StartedAt:  data.start.Format(utils.TimeFormat),
		DurationMs: data.end.Sub(data.start).Milliseconds(),
		Function:   funcStr,
		Plugin: Plugin{
			Id: pluginId,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func writeInferenceHistoryToPostgresDB(ctx context.Context, data inferenceHistory) error {
	return WithTx(ctx, func(tx pgx.Tx) error {
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

		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return err
		}

		return nil
	})
}

func QueryPlugins(ctx context.Context) ([]Plugin, error) {
	if globalModusDbEngine == nil {
		return nil, nil
	}
	_, plugins, err := modusgraph.Query[Plugin](ctx, globalModusDbEngine, modusgraph.QueryParams{})
	return plugins, err
}

func QueryInferences(ctx context.Context) ([]Inference, error) {
	if globalModusDbEngine == nil {
		return nil, nil
	}
	_, inferences, err := modusgraph.Query[Inference](ctx, globalModusDbEngine, modusgraph.QueryParams{})
	return inferences, err
}
