/*
 * Copyright 2023 Hypermode, Inc.
 */

//go:generate go run ./.tools/generate_version

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
		Msg("Starting Hypermode Runtime.")

	// Load environment variables from plugins path
	// note: .env file is optional, so don't log if it's not found
	err := godotenv.Load(filepath.Join(config.StoragePath, ".env"))
	if err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("Error reading .env file.  Ignoring.")
	}

	// Initialize Sentry
	rootSourcePath := getRootSourcePath()
	utils.InitSentry(rootSourcePath)
	defer utils.FlushSentryEvents()

	// Start the background services
	ctx := services.Start()
	defer services.Stop(ctx)

	// Set local mode if debugging is enabled
	local := utils.HypermodeDebugEnabled()

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
