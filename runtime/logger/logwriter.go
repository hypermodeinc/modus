/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package logger

import (
	"strings"

	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/rs/zerolog"
)

type logWriter struct {
	buffer *strings.Builder
	logger *zerolog.Logger
	level  zerolog.Level
}

// NewLogWriter creates a new log writer that writes to the given logger with the given level.
func NewLogWriter(logger *zerolog.Logger, level zerolog.Level) *logWriter {
	buffer := &strings.Builder{}
	return &logWriter{buffer, logger, level}
}

func (w logWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	for _, b := range p {
		if b == '\n' {
			line := w.buffer.String()
			if len(line) > 0 {
				w.logMessage(line)
				w.buffer.Reset()
			}
		} else {
			w.buffer.WriteByte(b)
		}
	}
	return
}

func (w logWriter) logMessage(line string) {
	l, message := utils.SplitConsoleOutputLine(line)
	level := ParseLevel(l)
	if level == zerolog.NoLevel {
		level = w.level
	}

	w.logger.
		WithLevel(level).
		Str("text", message).
		Msg("Message logged from function.")
}

func ParseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.NoLevel
	}
}
