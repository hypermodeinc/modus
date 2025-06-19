package db

import (
	"github.com/go-logr/logr"
	"github.com/rs/zerolog"
)

type Zerologr struct {
	logger zerolog.Logger
	level  int
}

// NewZerologr wraps a zerolog.Logger into a logr.Logger
func NewZerologr(l zerolog.Logger) logr.Logger {
	return logr.New(&Zerologr{logger: l, level: 0})
}

// Implement the logr.LogSink interface

func (z *Zerologr) Init(info logr.RuntimeInfo) {
	// No initialization needed
}

func (z *Zerologr) Enabled(level int) bool {
	return level <= z.level
}

func (z *Zerologr) Info(level int, msg string, keysAndValues ...interface{}) {
	if !z.Enabled(level) {
		return
	}
	evt := z.logger.Info()
	for i := 0; i < len(keysAndValues); i += 2 {
		key, val := keysAndValues[i], keysAndValues[i+1]
		evt.Interface(key.(string), val)
	}
	evt.Msg(msg)
}

func (z *Zerologr) Error(err error, msg string, keysAndValues ...interface{}) {
	evt := z.logger.Error().Err(err)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, val := keysAndValues[i], keysAndValues[i+1]
		evt.Interface(key.(string), val)
	}
	evt.Msg(msg)
}

func (z *Zerologr) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return &Zerologr{
		logger: z.logger.With().Fields(keysAndValues).Logger(),
		level:  z.level,
	}
}

func (z *Zerologr) WithName(name string) logr.LogSink {
	return &Zerologr{
		logger: z.logger.With().Str("logger", name).Logger(),
		level:  z.level,
	}
}

func (z *Zerologr) V(level int) logr.LogSink {
	return &Zerologr{
		logger: z.logger,
		level:  level,
	}
}
