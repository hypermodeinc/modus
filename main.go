/*
 * Copyright 2023 Hypermode, Inc.
 */

package main

import (
	"context"
	"os"
	"path"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/functions"
	"hmruntime/host"
	"hmruntime/logger"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	config.ParseCommandLineFlags()
	log := logger.Initialize()

	log.Info().Msg("Starting Hypermode Runtime.")

	// Validate configuration
	if config.PluginsPath == "" && config.S3Bucket == "" {
		log.Fatal().Msg("A plugins path and/or S3 bucket are required.  Exiting.")
	}

	if config.S3Bucket == "" {
		if _, err := os.Stat(config.PluginsPath); os.IsNotExist(err) {
			log.Info().
				Str("path", config.PluginsPath).
				Msg("Creating plugins directory.")
			err := os.MkdirAll(config.PluginsPath, 0755)
			if err != nil {
				log.Fatal().Err(err).
					Msg("Failed to create plugins directory.  Exiting.")
			}
		} else {
			log.Info().
				Str("path", config.PluginsPath).
				Msg("Found plugins directory.")
		}
	}

	// Load environment variables from plugins path
	// note: .env file is optional, so don't log if it's not found
	err := godotenv.Load(path.Join(config.PluginsPath, ".env"))
	if err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("Error reading .env file.  Ignoring.")
	}

	// Initialize the AWS configuration if we're using any AWS functionality
	if config.S3Bucket != "" || config.UseAwsSecrets {
		err = aws.Initialize(ctx)
		if err != nil {
			log.Info().Err(err).Msg("AWS functionality will be disabled.")
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

	// Watch for plugin changes
	err = host.WatchForPluginChanges(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to watch for plugin changes.  Exiting.")
	}

	// Watch for hypermode.json changes
	err = host.WatchForJsonChanges(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to watch for hypermode.json changes.  Exiting.")
	}

	// Start the web server
	err = startServer(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.  Exiting.")
	}

	// TODO: Shutdown gracefully
}
