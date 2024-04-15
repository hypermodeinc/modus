/*
 * Copyright 2024 Hypermode, Inc.
 */

package engine

import (
	"context"
	"hmruntime/logger"
	"unsafe"

	"github.com/jensneuse/abstractlogger"
	"github.com/rs/zerolog"
)

func NewLoggerAdapter(ctx context.Context) *LoggerAdapter {
	return &LoggerAdapter{
		l: logger.Get(ctx),
	}
}

type LoggerAdapter struct {
	l *zerolog.Logger
}

func (l *LoggerAdapter) Debug(msg string, fields ...abstractlogger.Field) {
	l.l.Debug().Fields(l.fields(fields)).Msg(msg)
}

func (l *LoggerAdapter) Warn(msg string, fields ...abstractlogger.Field) {
	l.l.Warn().Fields(l.fields(fields)).Msg(msg)
}

func (l *LoggerAdapter) Info(msg string, fields ...abstractlogger.Field) {
	l.l.Info().Fields(l.fields(fields)).Msg(msg)
}

func (l *LoggerAdapter) Error(msg string, fields ...abstractlogger.Field) {
	l.l.Error().Fields(l.fields(fields)).Msg(msg)
}

func (l *LoggerAdapter) Fatal(msg string, fields ...abstractlogger.Field) {
	l.l.Fatal().Fields(l.fields(fields)).Msg(msg)
}

func (l *LoggerAdapter) Panic(msg string, fields ...abstractlogger.Field) {
	l.l.Panic().Fields(l.fields(fields)).Msg(msg)
}

func (l *LoggerAdapter) LevelLogger(level abstractlogger.Level) abstractlogger.LevelLogger {
	return &LevelLoggerAdapter{
		l:     l.l,
		level: level,
	}
}

func (l *LoggerAdapter) fields(fields []abstractlogger.Field) map[string]interface{} {
	out := make(map[string]interface{}, len(fields))
	for _, f := range fields {

		lf := *convertLoggerField(&f)

		switch lf.kind {
		case abstractlogger.StringField:
			out[lf.key] = lf.stringValue
		case abstractlogger.ByteStringField:
			out[lf.key] = lf.byteValue
		case abstractlogger.IntField:
			out[lf.key] = lf.intValue
		case abstractlogger.BoolField:
			out[lf.key] = lf.intValue != 0
		case abstractlogger.ErrorField, abstractlogger.NamedErrorField:
			out[lf.key] = lf.errorValue
		case abstractlogger.StringsField:
			out[lf.key] = lf.stringsValue
		default:
			out[lf.key] = lf.interfaceValue
		}

	}
	return out
}

type LevelLoggerAdapter struct {
	l     *zerolog.Logger
	level abstractlogger.Level
}

func (s *LevelLoggerAdapter) Println(v ...interface{}) {
	switch s.level {
	case abstractlogger.DebugLevel:
		s.l.Debug().Msgf("%v", v[0])
	case abstractlogger.InfoLevel:
		s.l.Info().Msgf("%v", v[0])
	case abstractlogger.WarnLevel:
		s.l.Warn().Msgf("%v", v[0])
	case abstractlogger.ErrorLevel:
		s.l.Error().Msgf("%v", v[0])
	case abstractlogger.FatalLevel:
		s.l.Fatal().Msgf("%v", v[0])
	case abstractlogger.PanicLevel:
		s.l.Panic().Msgf("%v", v[0])
	}
}

func (s *LevelLoggerAdapter) Printf(format string, v ...interface{}) {
	switch s.level {
	case abstractlogger.DebugLevel:
		s.l.Debug().Msgf(format, v...)
	case abstractlogger.InfoLevel:
		s.l.Info().Msgf(format, v...)
	case abstractlogger.WarnLevel:
		s.l.Warn().Msgf(format, v...)
	case abstractlogger.ErrorLevel:
		s.l.Error().Msgf(format, v...)
	case abstractlogger.FatalLevel:
		s.l.Fatal().Msgf(format, v...)
	case abstractlogger.PanicLevel:
		s.l.Panic().Msgf(format, v...)
	}
}

// Some unsafe code to access the fields of the abstractlogger.Field struct.
// This is a necessary workaround because the struct's fields are not accessible.
// See https://github.com/jensneuse/abstractlogger/issues/2

type loggerField struct {
	kind           abstractlogger.FieldKind
	key            string
	stringValue    string
	stringsValue   []string
	intValue       int64
	byteValue      []byte
	interfaceValue interface{}
	errorValue     error
}

func convertLoggerField(f *abstractlogger.Field) *loggerField {
	p := unsafe.Pointer(f)
	return (*loggerField)(p)
}
