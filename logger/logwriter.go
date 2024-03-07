/*
 * Copyright 2024 Hypermode, Inc.
 */

package logger

import (
	"strings"

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
			s := w.buffer.String()
			if len(s) > 0 {
				w.logMessage(s)
				w.buffer.Reset()
			}
		} else {
			w.buffer.WriteByte(b)
		}
	}
	return
}

func (w logWriter) logMessage(s string) {

	// Check if the message is prefixed with a known level
	// and log it with the corresponding level if it does.
	a := strings.SplitAfterN(s, ": ", 2)
	if len(a) == 2 {
		switch a[0] {
		case "Debug: ":
			w.logger.Debug().Msg(a[1])
			return
		case "Info: ":
			w.logger.Info().Msg(a[1])
			return
		case "Warning: ":
			w.logger.Warn().Msg(a[1])
			return
		case "Error: ":
			w.logger.Error().Msg(a[1])
			return
		}
	}

	// Otherwise, log the message with the configured level.
	w.logger.WithLevel(w.level).CallerSkipFrame(1).Msg(s)
}
