/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/httpserver"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/services"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func main() {

	// Initialize the configuration
	config.Initialize()

	// Initialize the logger
	log := logger.Initialize()
	log.Info().
		Str("version", config.GetVersionNumber()).
		Str("environment", config.GetEnvironmentName()).
		Msg("Starting Modus Runtime.")

	// Load environment variables from .env file(s)
	config.LoadEnvFiles(log)

	// Initialize Sentry (if enabled)
	rootSourcePath := app.GetRootSourcePath()
	utils.InitSentry(rootSourcePath)
	defer utils.FlushSentryEvents()

	// Start the background services
	ctx := services.Start()
	defer services.Stop(ctx)

	// Set local mode in development
	local := config.IsDevEnvironment()

	// Start the HTTP server to listen for requests.
	// Note, this function blocks, and handles shutdown gracefully.
	httpserver.Start(ctx, local)
}
