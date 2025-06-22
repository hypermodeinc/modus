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
	"slices"
	"strings"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/rs/zerolog"
	actorLog "github.com/tochemey/goakt/v3/log"
)

// some messages are ignored during shutdown because they are expected
var shutdownIgnoredMessages = []string{
	"Failed to acquire semaphore: context canceled",
	" is down. modus is going to shutdown.",
}

func newActorLogger(logger *zerolog.Logger) *actorLogger {

	var minLevel zerolog.Level
	if utils.EnvVarFlagEnabled("MODUS_DEBUG_ACTORS") {
		minLevel = zerolog.DebugLevel
	} else {
		// goakt info level is too noisy, so default to show warnings and above
		minLevel = zerolog.WarnLevel
	}

	l := logger.Level(minLevel).With().Str("component", "actors").Logger()
	return &actorLogger{logger: &l}
}

type actorLogger struct {
	logger       *zerolog.Logger
	paused       bool
	shuttingDown bool
}

func (al *actorLogger) Pause() {
	al.paused = true
}

func (al *actorLogger) Resume() {
	al.paused = false
}

func (al *actorLogger) writeToLog(level zerolog.Level, msg string) {
	if al.paused {
		return
	}
	if al.shuttingDown && slices.ContainsFunc(shutdownIgnoredMessages, func(s string) bool {
		return strings.Contains(msg, s)
	}) {
		return
	}
	al.logger.WithLevel(level).Msg(msg)
}

func (al *actorLogger) Debug(v ...any) {
	al.writeToLog(zerolog.DebugLevel, fmt.Sprint(v...))
}

func (al *actorLogger) Debugf(format string, v ...any) {
	al.writeToLog(zerolog.DebugLevel, fmt.Sprintf(format, v...))
}

func (al *actorLogger) Info(v ...any) {
	al.writeToLog(zerolog.InfoLevel, fmt.Sprint(v...))
}

func (al *actorLogger) Infof(format string, v ...any) {
	al.writeToLog(zerolog.InfoLevel, fmt.Sprintf(format, v...))
}

func (al *actorLogger) Warn(v ...any) {
	al.writeToLog(zerolog.WarnLevel, fmt.Sprint(v...))
}

func (al *actorLogger) Warnf(format string, v ...any) {
	al.writeToLog(zerolog.WarnLevel, fmt.Sprintf(format, v...))
}

func (al *actorLogger) Error(v ...any) {
	al.writeToLog(zerolog.ErrorLevel, fmt.Sprint(v...))
}

func (al *actorLogger) Errorf(format string, v ...any) {
	al.writeToLog(zerolog.ErrorLevel, fmt.Sprintf(format, v...))
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
