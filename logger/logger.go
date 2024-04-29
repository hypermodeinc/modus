/*
 * Copyright 2024 Hypermode, Inc.
 */

package logger

import (
	"context"
	"io"
	"os"
	"time"

	"hmruntime/config"
	"hmruntime/plugins"
	"hmruntime/utils"

	zls "github.com/archdx/zerolog-sentry"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// these are the field names that will be used in the logs
const executionIdKey = "execution_id"
const pluginNameKey = "plugin"
const buildIdKey = "build_id"

var zlsCloser io.Closer

func Initialize() *zerolog.Logger {
	var writer io.Writer
	if config.UseJsonLogging {
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
		writer = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "2006-01-02 15:04:05.000 -07:00",
		}
	}

	// Use zerolog-sentry to route error, fatal, and panic logs to Sentry.
	zlsWriter, err := zls.NewWithHub(sentry.CurrentHub(),
		zls.WithEnvironment(config.GetEnvironmentName()),
		zls.WithRelease(config.GetVersionNumber()))
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

func Get(ctx context.Context) *zerolog.Logger {
	executionId, eidFound := ctx.Value(utils.ExecutionIdContextKey).(string)
	plugin, pluginFound := ctx.Value(utils.PluginContextKey).(*plugins.Plugin)

	// If no context values, just return the global logger.
	if !eidFound && !pluginFound {
		return &log.Logger
	}

	// Create a logger context with the context values.
	lc := log.Logger.With()

	if executionId != "" {
		lc = lc.Str(executionIdKey, executionId)
	}

	var buildId string
	var pluginName string
	if pluginFound {
		buildId = plugin.BuildId()
		pluginName = plugin.Name()
	}

	if buildId != "" {
		lc = lc.Str(buildIdKey, buildId)
	}

	if pluginName != "" {
		lc = lc.Str(pluginNameKey, pluginName)
	}

	// Return a logger built from the context.
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
