/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/envfiles"
	"github.com/hypermodeinc/modus/runtime/httpserver"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/services"
)

func main() {

	// Initialize the app configuration (command-line flags, etc.)
	app.Initialize()

	// Create the main background context
	ctx := context.Background()

	// Initialize the logger
	log := logger.Initialize()
	log.Info().
		Str("version", app.VersionNumber()).
		Str("environment", app.Config().Environment()).
		Msg("Starting Modus Runtime.")

	err := envfiles.LoadEnvFiles(ctx)
	if err != nil {
		logger.Warn(ctx, err).Msg("Failed to load environment files.")
	}

	// Initialize Sentry (if enabled)
	sentryutils.InitializeSentry()
	defer sentryutils.CloseSentry()

	// Get the main handler for the HTTP server before starting the services,
	// so it can register the endpoints as the manifest is loaded.
	mux := httpserver.GetMainHandler()

	// Start the background services
	ctx = services.Start(ctx)
	defer services.Stop(ctx)

	// Set local mode in development
	local := app.IsDevEnvironment()

	// Start the HTTP server to listen for requests.
	// Note, this function blocks, and handles shutdown gracefully.
	httpserver.Start(ctx, mux, local)
}
