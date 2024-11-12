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
	"context"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/envfiles"
	"github.com/hypermodeinc/modus/runtime/httpserver"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/services"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func main() {

	// Initialize the configuration
	config.Initialize()

	// Create the main background context
	ctx := context.Background()

	// Initialize the logger
	log := logger.Initialize()
	log.Info().
		Str("version", config.GetVersionNumber()).
		Str("environment", config.GetEnvironmentName()).
		Msg("Starting Modus Runtime.")

	err := envfiles.LoadEnvFiles(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load environment files.")
	}

	// Initialize Sentry (if enabled)
	rootSourcePath := app.GetRootSourcePath()
	utils.InitSentry(rootSourcePath)
	defer utils.FlushSentryEvents()

	// Start the background services
	ctx = services.Start(ctx)
	defer services.Stop(ctx)

	// Set local mode in development
	local := config.IsDevEnvironment()

	// Start the HTTP server to listen for requests.
	// Note, this function blocks, and handles shutdown gracefully.
	httpserver.Start(ctx, local)
}
