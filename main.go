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
	"hmruntime/manifestdata"
	"hmruntime/secrets"
	"hmruntime/server"
	"hmruntime/storage"
	"hmruntime/utils"
	"hmruntime/vector"
	"hmruntime/wasmhost"

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

	// Initialize the secrets provider
	secrets.Initialize(ctx)

	// Initialize the storage provider
	storage.Initialize(ctx)

	// Initialize the metadata database
	db.Initialize(ctx)

	// Watch for function registration requests
	functions.MonitorRegistration(ctx)

	// Load app data and monitor for changes
	manifestdata.MonitorManifestFile(ctx)

	// Load plugins and monitor for changes
	wasmhost.MonitorPlugins(ctx)

	// Initialize the GraphQL engine
	graphql.Initialize()

	// Initialize in mem vector factory
	vector.InitializeIndexFactory()
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
