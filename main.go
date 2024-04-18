/*
 * Copyright 2023 Hypermode, Inc.
 */

package main

import (
	"context"
	"os"
	"path/filepath"

	"hmruntime/appdata"
	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/functions"
	"hmruntime/graphql"
	"hmruntime/host"
	"hmruntime/logger"
	"hmruntime/server"
	"hmruntime/storage"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	config.ParseCommandLineFlags()
	log := logger.Initialize()

	log.Info().Msg("Starting Hypermode Runtime.")

	// Load environment variables from plugins path
	// note: .env file is optional, so don't log if it's not found
	err := godotenv.Load(filepath.Join(config.StoragePath, ".env"))
	if err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("Error reading .env file.  Ignoring.")
	}

	// Initialize the AWS configuration if we're using any AWS functionality
	if config.UseAwsStorage || config.UseAwsSecrets {
		err = aws.Initialize(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize AWS.  Exiting.")
		}
	}

	// Initialize the WebAssembly runtime
	err = host.InitWasmRuntime(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize the WebAssembly runtime.  Exiting.")
	}
	defer host.WasmRuntime.Close(ctx)

	// Connect Hypermode host functions
	err = functions.InstantiateHostFunctions(ctx, host.WasmRuntime)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to instantiate host functions.  Exiting.")
	}

	// Initialize the storage system
	storage.Initialize()

	// Watch for function registration requests
	functions.MonitorRegistration(ctx)

	// Load app data and monitor for changes
	appdata.MonitorAppDataFiles(ctx)

	// Load plugins and monitor for changes
	host.MonitorPlugins(ctx)

	// Initialize the GraphQL engine
	graphql.Initialize()

	// Start the web server
	err = server.Start(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.  Exiting.")
	}

	// TODO: Shutdown gracefully
}
