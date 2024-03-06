/*
 * Copyright 2024 Hypermode, Inc.
 */

package logger

import (
	"context"
	"hmruntime/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type contextKey string

const correlationIdKey = "correlation_id"
const CorrelationIdContextKey contextKey = correlationIdKey

func Initialize() *zerolog.Logger {
	if !config.UseJsonLogging {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	}

	return &log.Logger
}

func Get(ctx context.Context) *zerolog.Logger {
	correlationId, ok := ctx.Value(CorrelationIdContextKey).(string)
	if ok && correlationId != "" {
		l := log.Logger.With().Str(correlationIdKey, correlationId).Logger()
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
