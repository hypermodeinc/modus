/*
 * Copyright 2023 Hypermode, Inc.
 */

package main

import (
	"context"
	"os"
	"path/filepath"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/functions"
	"hmruntime/host"
	"hmruntime/logger"
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

	// Initialize the storage system
	storage.Initialize()

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

	// Load json
	err = host.LoadJsons(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load hypermode.json.  Exiting.")
	}

	// Load plugins
	err = host.LoadPlugins(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load plugins.  Exiting.")
	}

	// Watch for registration requests
	functions.MonitorRegistration(ctx)

	// Watch for schema changes
	functions.MonitorGqlSchema(ctx)

	// // Watch for plugin changes
	// err = host.WatchForPluginChanges(ctx)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to watch for plugin changes.  Exiting.")
	// }

	// // Watch for hypermode.json changes
	// err = host.WatchForJsonChanges(ctx)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to watch for hypermode.json changes.  Exiting.")
	// }

	// Start the web server
	err = startServer(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.  Exiting.")
	}

	// TODO: Shutdown gracefully
}
