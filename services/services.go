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
	"hmruntime/sqlclient"
	"hmruntime/storage"
	"hmruntime/utils"
	"hmruntime/wasmhost"
)

// Starts any services that need to be started when the runtime starts.
func Start(ctx context.Context) {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// None of these should block. If they need to do background work, they should start a goroutine internally.

	// NOTE: Initialization order is important in some cases.
	// If you need to change the order or add new services, be sure to test thoroughly.
	// Generally, new services should be added to the end of the list, unless there is a specific reason to do otherwise.

	wasmhost.InitWasmHost(ctx)
	hostfunctions.RegisterHostFunctions(ctx)
	aws.Initialize(ctx)
	secrets.Initialize(ctx)
	storage.Initialize(ctx)
	db.Initialize(ctx)
	collections.InitializeIndexFactory(ctx)
	manifestdata.MonitorManifestFile(ctx)
	pluginmanager.MonitorPlugins(ctx)
	graphql.Initialize()
}

// Stops any services that need to be stopped when the runtime stops.
func Stop(ctx context.Context) {

	// NOTE: Stopping services also has an order dependency.
	// If you need to change the order or add new services, be sure to test thoroughly.

	// Unlike start, these should each block until they are fully stopped.

	collections.CloseIndexFactory(ctx)
	sqlclient.ShutdownPGPools()
	wasmhost.RuntimeInstance.Close(ctx)
	logger.Close()
	db.Stop(ctx)
}
