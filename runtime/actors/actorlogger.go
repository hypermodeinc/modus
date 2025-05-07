/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package actors

import (
	"fmt"
	"io"
	"log"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/rs/zerolog"
	actorLog "github.com/tochemey/goakt/v3/log"
)

func newActorLogger(logger *zerolog.Logger) *actorLogger {

	var minLevel zerolog.Level
	if utils.EnvVarFlagEnabled("MODUS_DEBUG_ACTORS") {
		minLevel = zerolog.DebugLevel
	} else {
		// goakt info level is too noisy, so default to show only errors
		minLevel = zerolog.ErrorLevel
	}

	l := logger.Level(minLevel).With().Str("component", "actors").Logger()
	return &actorLogger{&l}
}

type actorLogger struct {
	logger *zerolog.Logger
}

func (al *actorLogger) Debug(v ...any) {
	al.logger.Debug().Msg(fmt.Sprint(v...))
}

func (al *actorLogger) Debugf(format string, v ...any) {
	al.logger.Debug().Msgf(format, v...)
}

func (al *actorLogger) Info(v ...any) {
	al.logger.Info().Msg(fmt.Sprint(v...))
}

func (al *actorLogger) Infof(format string, v ...any) {
	al.logger.Info().Msgf(format, v...)
}

func (al *actorLogger) Warn(v ...any) {
	al.logger.Warn().Msg(fmt.Sprint(v...))
}

func (al *actorLogger) Warnf(format string, v ...any) {
	al.logger.Warn().Msgf(format, v...)
}

func (al *actorLogger) Error(v ...any) {
	al.logger.Error().Msg(fmt.Sprint(v...))
}

func (al *actorLogger) Errorf(format string, v ...any) {
	al.logger.Error().Msgf(format, v...)
}

func (al *actorLogger) Fatal(v ...any) {
	al.logger.Fatal().Msg(fmt.Sprint(v...))
}

func (al *actorLogger) Fatalf(format string, v ...any) {
	al.logger.Fatal().Msgf(format, v...)
}

func (al *actorLogger) Panic(v ...any) {
	al.logger.Panic().Msg(fmt.Sprint(v...))
}

func (al *actorLogger) Panicf(format string, v ...any) {
	al.logger.Panic().Msgf(format, v...)
}

func (al *actorLogger) LogLevel() actorLog.Level {
	switch al.logger.GetLevel() {
	case zerolog.DebugLevel:
		return actorLog.DebugLevel
	case zerolog.InfoLevel:
		return actorLog.InfoLevel
	case zerolog.WarnLevel:
		return actorLog.WarningLevel
	case zerolog.ErrorLevel:
		return actorLog.ErrorLevel
	case zerolog.FatalLevel:
		return actorLog.FatalLevel
	case zerolog.PanicLevel:
		return actorLog.PanicLevel
	default:
		return actorLog.InvalidLevel
	}
}

func (al *actorLogger) LogOutput() []io.Writer {
	w := logger.NewLogWriter(al.logger, al.logger.GetLevel())
	return []io.Writer{w}
}

func (al *actorLogger) StdLogger() *log.Logger {
	w := logger.NewLogWriter(al.logger, al.logger.GetLevel())
	return log.New(w, "", 0)
}
