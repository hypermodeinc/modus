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
	"os"
	"sync"
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Initialize() *zerolog.Logger {
	if app.Config().UseJsonLogging() {
		// In JSON mode, we'll log UTC with millisecond precision.
		// Note that Go uses this specific value for its formatting exemplars.
		zerolog.TimeFieldFormat = utils.TimeFormat
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().UTC()
		}
		log.Logger = log.Logger.Output(os.Stderr)
	} else {
		// In console mode, we can use local time and be a bit prettier.
		// We'll still log with millisecond precision.
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		if app.IsDevEnvironment() {
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
				"agent",
				"agent_id",
				"execution_id",
				"duration_ms",
				"cluster_mode",
				"cluster_host",
				"discovery_port",
				"remoting_port",
				"peers_port",
			}
		} else {
			consoleWriter.TimeFormat = "2006-01-02 15:04:05.000 -07:00"
		}

		log.Logger = log.Logger.Output(consoleWriter)
	}

	// Log the runtime version to every log line, except in development.
	if !app.IsDevEnvironment() {
		log.Logger = log.Logger.With().
			Str("runtime_version", app.VersionNumber()).
			Logger()
	}

	return &log.Logger
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

func Warn(ctx context.Context, errs ...error) *zerolog.Event {
	switch len(errs) {
	case 0:
		return Get(ctx).Warn()
	case 1:
		err := errs[0]
		if err == nil {
			return Get(ctx).Warn()
		}
		return Get(ctx).Warn().Err(err)
	default:
		return Get(ctx).Warn().Errs("errors", errs)
	}
}

func Error(ctx context.Context, errs ...error) *zerolog.Event {
	switch len(errs) {
	case 0:
		return Get(ctx).Error()
	case 1:
		err := errs[0]
		if err == nil {
			return Get(ctx).Error()
		}
		return Get(ctx).Err(err)
	default:
		return Get(ctx).Error().Errs("errors", errs)
	}
}

func Fatal(ctx context.Context, errs ...error) *zerolog.Event {
	switch len(errs) {
	case 0:
		return Get(ctx).Fatal()
	case 1:
		err := errs[0]
		if err == nil {
			return Get(ctx).Fatal()
		}
		return Get(ctx).Fatal().Err(err)
	default:
		return Get(ctx).Fatal().Errs("errors", errs)
	}
}

func Debugf(msg string, v ...any) {
	log.Debug().Msgf(msg, v...)
}

func Infof(msg string, v ...any) {
	log.Info().Msgf(msg, v...)
}

func Warnf(msg string, v ...any) {
	log.Warn().Msgf(msg, v...)
}

func Errorf(msg string, v ...any) {
	log.Error().Msgf(msg, v...)
}

func Fatalf(msg string, v ...any) {
	log.Fatal().Msgf(msg, v...)
}

func Errf(err error, msg string, v ...any) {
	log.Err(err).Msgf(msg, v...)
}
