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
	"hmruntime/collections"
	"hmruntime/config"
	"hmruntime/db"
	"hmruntime/graphql"
	"hmruntime/hostfunctions"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/pluginmanager"
	"hmruntime/secrets"
	"hmruntime/server"
	"hmruntime/storage"
	"hmruntime/utils"
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

	// Set local mode if debugging is enabled
	local := utils.HypermodeDebugEnabled()

	// Start the HTTP server to listen for requests.
	// Note, this function blocks, and handles shutdown gracefully.
	server.Start(ctx, local)
}

func initRuntimeServices(ctx context.Context) {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// Initialize the WebAssembly runtime
	wasmhost.InitWasmHost(ctx)

	// Register the host functions with the runtime
	hostfunctions.RegisterHostFunctions(ctx)

	// Initialize AWS functionality
	aws.Initialize(ctx)

	// Initialize the secrets provider
	secrets.Initialize(ctx)

	// Initialize the storage provider
	storage.Initialize(ctx)

	// Initialize the metadata database
	db.Initialize(ctx)

	// Initialize in mem vector factory
	collections.InitializeIndexFactory(ctx)

	// Load app data and monitor for changes
	onChange := func([]error) {
		hostfunctions.ShutdownPools()
	}
	manifestdata.MonitorManifestFile(ctx, onChange)

	// Load plugins and monitor for changes
	pluginmanager.MonitorPlugins(ctx)

	// Initialize the GraphQL engine
	graphql.Initialize()
}

func stopRuntimeServices(ctx context.Context) {
	collections.CloseIndexFactory(ctx)
	wasmhost.RuntimeInstance.Close(ctx)
	logger.Close()

	db.Stop()
	hostfunctions.ShutdownPools()
}

func getRootSourcePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return path.Dir(filename) + "/"
}
