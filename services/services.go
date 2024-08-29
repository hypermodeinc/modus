/*
 * Copyright 2024 Hypermode, Inc.
 */

package services

import (
	"context"

	"hypruntime/aws"
	"hypruntime/collections"
	"hypruntime/db"
	"hypruntime/dgraphclient"
	"hypruntime/graphql"
	"hypruntime/hostfunctions"
	"hypruntime/logger"
	"hypruntime/manifestdata"
	"hypruntime/pluginmanager"
	"hypruntime/secrets"
	"hypruntime/sqlclient"
	"hypruntime/storage"
	"hypruntime/utils"
	"hypruntime/wasmhost"
)

// Starts any services that need to be started when the runtime starts.
func Start(ctx context.Context) {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// Init the wasm host before anything else.
	registrations := hostfunctions.GetRegistrations()
	wasmhost.InitWasmHost(ctx, registrations...)

	// None of these should block. If they need to do background work, they should start a goroutine internally.

	// NOTE: Initialization order is important in some cases.
	// If you need to change the order or add new services, be sure to test thoroughly.
	// Generally, new services should be added to the end of the list, unless there is a specific reason to do otherwise.

	sqlclient.Initialize()
	dgraphclient.Initialize()
	aws.Initialize(ctx)
	secrets.Initialize(ctx)
	storage.Initialize(ctx)
	db.Initialize(ctx)
	collections.Initialize(ctx)
	manifestdata.MonitorManifestFile(ctx)
	pluginmanager.Initialize(ctx)
	graphql.Initialize()
}

// Stops any services that need to be stopped when the runtime stops.
func Stop(ctx context.Context) {

	// NOTE: Stopping services also has an order dependency.
	// If you need to change the order or add new services, be sure to test thoroughly.

	// Unlike start, these should each block until they are fully stopped.

	collections.Shutdown(ctx)
	sqlclient.ShutdownPGPools()
	dgraphclient.ShutdownConns()
	wasmhost.GlobalWasmHost.Close(ctx)
	logger.Close()
	db.Stop(ctx)
}
