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

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/db"
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

	// Initialize Sentry
	rootSourcePath := getRootSourcePath()
	utils.InitSentry(rootSourcePath)
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
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// Initialize the WebAssembly runtime
	wasmhost.InitWasmHost(ctx)

	// Initialize AWS functionality
	aws.Initialize(ctx)

	// Initialize the inference history database
	db.Initialize(ctx)

	// Initialize the storage system
	storage.Initialize(ctx)

	// Watch for function registration requests
	functions.MonitorRegistration()

	// Load app data and monitor for changes
	manifest.MonitorAppDataFiles()

	// Load plugins and monitor for changes
	wasmhost.MonitorPlugins()

	// Initialize the GraphQL engine
	graphql.Initialize()
}

func stopRuntimeServices(ctx context.Context) {
	wasmhost.RuntimeInstance.Close(ctx)
	logger.Close()

	db.Stop()
}

func getRootSourcePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return path.Dir(filename) + "/"
}
