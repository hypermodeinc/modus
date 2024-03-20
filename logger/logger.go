/*
 * Copyright 2024 Hypermode, Inc.
 */

package logger

import (
	"context"
	"os"
	"time"

	"hmruntime/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type contextKey string

const executionIdKey = "execution_id"
const pluginNameKey = "plugin"
const buildIdKey = "build_id"

const ExecutionIdContextKey contextKey = executionIdKey
const PluginNameContextKey contextKey = pluginNameKey
const BuildIdContextKey contextKey = buildIdKey

func Initialize() *zerolog.Logger {
	if config.UseJsonLogging {
		// In JSON mode, we'll log UTC with millisecond precision.
		// Note that Go uses this specific value for its formatting exemplars.
		zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z"
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().UTC()
		}
	} else {
		// In console mode, we can use local time and be a bit prettier.
		// We'll still log with millisecond precision.
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "2006-01-02 15:04:05.999 -07:00",
		})
	}

	return &log.Logger
}

func Get(ctx context.Context) *zerolog.Logger {
	executionId, found1 := ctx.Value(ExecutionIdContextKey).(string)
	buildId, found2 := ctx.Value(BuildIdContextKey).(string)
	pluginName, found3 := ctx.Value(PluginNameContextKey).(string)

	// If no context values, just return the global logger.
	if !found1 && !found2 && !found3 {
		return &log.Logger
	}

	// Create a logger context with the context values.
	lc := log.Logger.With()

	if executionId != "" {
		lc = lc.Str(executionIdKey, executionId)
	}

	if buildId != "" {
		lc = lc.Str(buildIdKey, buildId)
	}

	if pluginName != "" {
		lc = lc.Str(pluginNameKey, pluginName)
	}

	// Return a logger built from the context.
	l := lc.Logger()
	return &l
}

func Trace(ctx context.Context) *zerolog.Event {
	return Get(ctx).Trace()
}

func Debug(ctx context.Context) *zerolog.Event {
	return Get(ctx).Debug()
}

func Info(ctx context.Context) *zerolog.Event {
	return Get(ctx).Info()
}

func Warn(ctx context.Context) *zerolog.Event {
	return Get(ctx).Warn()
}

func Error(ctx context.Context) *zerolog.Event {
	return Get(ctx).Error()
}

func Err(ctx context.Context, err error) *zerolog.Event {
	return Get(ctx).Err(err)
}

func Fatal(ctx context.Context) *zerolog.Event {
	return Get(ctx).Fatal()
}
