/*
 * Copyright 2024 Hypermode, Inc.
 */

package services

import (
	"context"

	"hmruntime/aws"
	"hmruntime/collections"
	"hmruntime/db"
	"hmruntime/graphql"
	"hmruntime/hostfunctions"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/pluginmanager"
	"hmruntime/secrets"
	"hmruntime/storage"
	"hmruntime/utils"
	"hmruntime/wasmhost"
)

func Start(ctx context.Context) {
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
	manifestdata.MonitorManifestFile(ctx)

	// Load plugins and monitor for changes
	pluginmanager.MonitorPlugins(ctx)

	// Initialize the GraphQL engine
	graphql.Initialize()
}

func Stop(ctx context.Context) {
	collections.CloseIndexFactory(ctx)
	wasmhost.RuntimeInstance.Close(ctx)
	logger.Close()

	db.Stop()
}
