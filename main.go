/*
 * Copyright 2023 Hypermode, Inc.
 */

package main

import (
	"context"
	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/functions"
	"hmruntime/host"
	"hmruntime/logger"
	"hmruntime/plugins"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	config.ParseCommandLineFlags()
	log := logger.Initialize()

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

	// Initialize the AWS configuration
	err := aws.Initialize(ctx)
	if err != nil {
		log.Info().Err(err).Msg("AWS functionality will be disabled.")
	}

	// Initialize the WebAssembly runtime
	host.WasmRuntime, err = plugins.InitWasmRuntime(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize the WebAssembly runtime.  Exiting.")
	}
	defer host.WasmRuntime.Close(ctx)

	if !aws.Enabled() {
		// Load environment variables from plugins path
		err = godotenv.Load(config.PluginsPath + "/.env")
		if err != nil {
			log.Info().Err(err).Msg("No .env file found.")
		}
	}

	// Load json
	err = plugins.LoadJsons(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load hypermode.json.  Exiting.")
	}

	// Load plugins
	err = plugins.LoadPlugins(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load plugins.  Exiting.")
	}

	// Watch for registration requests
	functions.MonitorRegistration(ctx)

	// Watch for schema changes
	functions.MonitorGqlSchema(ctx)

	// Watch for plugin changes
	err = plugins.WatchForPluginChanges(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to watch for plugin changes.  Exiting.")
	}

	// Watch for hypermode.json changes
	err = plugins.WatchForJsonChanges(ctx)
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
