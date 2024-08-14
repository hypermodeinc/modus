/*
 * Copyright 2024 Hypermode, Inc.
 */

package pluginmanager

import (
	"context"

	"hmruntime/logger"
	"hmruntime/plugins"

	"github.com/rs/zerolog"
)

func Initialize(ctx context.Context) {
	configureLogger()
	monitorPlugins(ctx)
}

func configureLogger() {
	logger.AddAdapter(func(ctx context.Context, lc zerolog.Context) zerolog.Context {

		if plugin := plugins.GetPlugin(ctx); plugin != nil {
			if buildId := plugin.BuildId(); buildId != "" {
				lc = lc.Str("build_id", buildId)
			}

			if pluginName := plugin.Name(); pluginName != "" {
				lc = lc.Str("plugin", pluginName)
			}
		}

		return lc
	})
}
