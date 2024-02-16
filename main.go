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
	"hmruntime/plugins"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})

	ctx := context.Background()
	config.ParseCommandLineFlags()

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

	// Initialize .env if in dev mode
	if os.Getenv("ENV") == "dev" {
		// Try to load the .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal().Err(err).Msg("Error loading .env file")
		}
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

	// Start the web server
	err = startServer()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.  Exiting.")
	}

	// TODO: Shutdown gracefully
}
