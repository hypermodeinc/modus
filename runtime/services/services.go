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
func Start() context.Context {
	// Create the main background context
	ctx := context.Background()

	// Note, we cannot start a Sentry transaction here, or it will also be used for the background services, post-initiation.

	// Init the wasm host and put it in context
	registrations := hostfunctions.GetRegistrations()
	host := wasmhost.InitWasmHost(ctx, registrations...)
	ctx = context.WithValue(ctx, utils.WasmHostContextKey, host)

	// Start the background services.  None of these should block. If they need to do background work, they should start a goroutine internally.

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

	return ctx
}

// Stops any services that need to be stopped when the runtime stops.
func Stop(ctx context.Context) {

	// Stop the wasm host first
	wasmhost.GetWasmHost(ctx).Close(ctx)

	// Stop the rest of the background services.
	// NOTE: Stopping services also has an order dependency.
	// If you need to change the order or add new services, be sure to test thoroughly.
	// Unlike start, these should each block until they are fully stopped.

	collections.Shutdown(ctx)
	sqlclient.ShutdownPGPools()
	dgraphclient.ShutdownConns()
	logger.Close()
	db.Stop(ctx)
}
