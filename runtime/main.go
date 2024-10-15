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
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/httpserver"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/middleware"
	"github.com/hypermodeinc/modus/runtime/services"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
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
	loadEnvFiles(log)

	// Initialize Sentry
	rootSourcePath := getRootSourcePath()
	utils.InitSentry(rootSourcePath)
	defer utils.FlushSentryEvents()

	// Start the background services
	ctx := services.Start()
	defer services.Stop(ctx)

	// Retrieve auth private keys
	middleware.Init(ctx)

	// Set local mode in development
	local := config.IsDevEnvironment()

	// Start the HTTP server to listen for requests.
	// Note, this function blocks, and handles shutdown gracefully.
	httpserver.Start(ctx, local)
}

func loadEnvFiles(log *zerolog.Logger) {
	envName := config.GetEnvironmentName()

	files := []string{
		".env." + envName + ".local",
		".env." + envName,
		".env.local",
		".env",
	}

	for _, file := range files {
		path := filepath.Join(config.AppPath, file)
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				log.Warn().Err(err).Msgf("Failed to load %s file.", file)
			}
		}
	}
}

func getRootSourcePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return path.Dir(filename) + "/"
}
