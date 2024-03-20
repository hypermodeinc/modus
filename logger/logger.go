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
const ExecutionIdContextKey contextKey = executionIdKey

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
	executionId, ok := ctx.Value(ExecutionIdContextKey).(string)
	if ok && executionId != "" {
		l := log.Logger.With().Str(executionIdKey, executionId).Logger()
		return &l
	}

	return &log.Logger
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
