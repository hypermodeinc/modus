/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package pluginmanager

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/plugins"

	"github.com/rs/zerolog"
)

func Initialize(ctx context.Context) {
	configureLogger()
	monitorPlugins(ctx)
}

func configureLogger() {
	logger.AddAdapter(func(ctx context.Context, lc zerolog.Context) zerolog.Context {

		if plugin, ok := plugins.GetPluginFromContext(ctx); ok {
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
