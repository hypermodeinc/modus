/*
 * Copyright 2023 Hypermode, Inc.
 */

//go:generate go run ./.tools/generate_version

package main

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"hmruntime/config"
	"hmruntime/logger"
	"hmruntime/server"
	"hmruntime/services"
	"hmruntime/utils"

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

	// Create the main background context
	ctx := context.Background()

	// Start the runtime services
	services.Start(ctx)
	defer services.Stop(ctx)

	// Start the HTTP server to listen for requests.
	// Note, this function blocks, and handles shutdown gracefully.
	server.Start(ctx)
}

func getRootSourcePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return path.Dir(filename) + "/"
}
