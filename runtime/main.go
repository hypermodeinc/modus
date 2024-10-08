/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

//go:generate go run ./tools/generate_version

package main

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/httpserver"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/services"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/joho/godotenv"
)

func main() {

	// Initialize the configuration
	config.Initialize()

	// Initialize the logger
	log := logger.Initialize()
	log.Info().
		Str("version", config.GetVersionNumber()).
		Msg("Starting Modus Runtime.")

	// Load environment variables from plugins path
	// note: .env file is optional, so don't log if it's not found
	err := godotenv.Load(filepath.Join(config.StoragePath, ".env"))
	if err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("Error reading .env file.  Ignoring.")
	}
	if config.IsDevEnvironment() {
		err = godotenv.Load(filepath.Join(config.StoragePath, ".env.local"))
		if err != nil && !os.IsNotExist(err) {
			log.Warn().Err(err).Msg("Error reading .env.local file.  Ignoring.")
		}
	}

	// Initialize Sentry
	rootSourcePath := getRootSourcePath()
	utils.InitSentry(rootSourcePath)
	defer utils.FlushSentryEvents()

	// Start the background services
	ctx := services.Start()
	defer services.Stop(ctx)

	// Set local mode if debugging is enabled
	local := utils.DebugModeEnabled()

	// Start the HTTP server to listen for requests.
	// Note, this function blocks, and handles shutdown gracefully.
	httpserver.Start(ctx, local)
}

func getRootSourcePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return path.Dir(filename) + "/"
}
