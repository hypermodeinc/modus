/*
 * Copyright 2023 Hypermode, Inc.
 */

//go:generate go run ./.tools/generate_version

package main

import (
	"context"
	"os"
	"path/filepath"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/functions"
	"hmruntime/graphql"
	"hmruntime/logger"
	"hmruntime/manifest"
	"hmruntime/server"
	"hmruntime/storage"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize Sentry
	utils.InitSentry()
	defer utils.FlushSentryEvents()

	// Initialize the runtime services
	ctx := context.Background()
	initRuntimeServices(ctx)
	defer stopRuntimeServices(ctx)

	// Start the HTTP server to listen for requests.
	// Note, this function blocks, and handles shutdown gracefully.
	server.Start(ctx)
}

func initRuntimeServices(ctx context.Context) {

	// Parse the command line flags
	config.ParseCommandLineFlags()

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

	// Instrument the rest of the startup process
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// Initialize the AWS configuration if we're using any AWS functionality
	if config.UseAwsStorage || config.UseAwsSecrets {
		err = aws.Initialize(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize AWS.  Exiting.")
		}
	}

	// Initialize the WebAssembly runtime
	wasmhost.InitWasmHost(ctx)

	// Initialize the storage system
	storage.Initialize(ctx)

	// Watch for function registration requests
	functions.MonitorRegistration(ctx)

	// Load app data and monitor for changes
	manifest.MonitorAppDataFiles(ctx)

	// Load plugins and monitor for changes
	wasmhost.MonitorPlugins(ctx)

	// Initialize the GraphQL engine
	graphql.Initialize(ctx)
}

func stopRuntimeServices(ctx context.Context) {
	wasmhost.RuntimeInstance.Close(ctx)
	logger.Close()
}
