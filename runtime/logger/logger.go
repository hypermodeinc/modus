/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package logger

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/utils"

	zls "github.com/archdx/zerolog-sentry"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var zlsCloser io.Closer

func Initialize() *zerolog.Logger {
	var writer io.Writer
	if app.Config().UseJsonLogging() {
		// In JSON mode, we'll log UTC with millisecond precision.
		// Note that Go uses this specific value for its formatting exemplars.
		zerolog.TimeFieldFormat = utils.TimeFormat
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().UTC()
		}
		writer = os.Stderr
	} else {
		// In console mode, we can use local time and be a bit prettier.
		// We'll still log with millisecond precision.
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		if app.Config().IsDevEnvironment() {
			consoleWriter.TimeFormat = "15:04:05.000"
			consoleWriter.FieldsExclude = []string{
				"build_id",
				"build_ts",
				"git_commit",
				"git_repo",
				"plugin",
				"user_visible",
			}
			consoleWriter.FieldsOrder = []string{
				"detail",
				"function",
				"execution_id",
				"duration_ms",
			}
		} else {
			consoleWriter.TimeFormat = "2006-01-02 15:04:05.000 -07:00"
		}

		writer = consoleWriter
	}

	// Log the runtime version to every log line, except in development.
	if !app.Config().IsDevEnvironment() {
		log.Logger = log.Logger.With().
			Str("runtime_version", app.VersionNumber()).
			Logger()
	}

	// Use zerolog-sentry to route error, fatal, and panic logs to Sentry.
	zlsWriter, err := zls.NewWithHub(sentry.CurrentHub(), zls.WithBreadcrumbs())
	if err != nil {
		logger := log.Logger.Output(writer)
		logger.Fatal().Err(err).Msg("Failed to initialize Sentry logger.")
	}
	zlsCloser = zlsWriter // so we can close it later, which flushes Sentry events
	log.Logger = log.Logger.Output(zerolog.MultiLevelWriter(writer, zlsWriter))

	return &log.Logger
}

func Close() {
	if zlsCloser != nil {
		zlsCloser.Close()
	}
}

var adapters []func(context.Context, zerolog.Context) zerolog.Context
var mu sync.RWMutex

func AddAdapter(adapter func(context.Context, zerolog.Context) zerolog.Context) {
	mu.Lock()
	defer mu.Unlock()
	adapters = append(adapters, adapter)
}

func Get(ctx context.Context) *zerolog.Logger {
	mu.RLock()
	defer mu.RUnlock()

	if len(adapters) == 0 {
		return &log.Logger
	}

	lc := log.Logger.With()
	for _, adapter := range adapters {
		lc = adapter(ctx, lc)
	}

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
