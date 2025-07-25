/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package services

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/actors"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/dgraphclient"
	"github.com/hypermodeinc/modus/runtime/envfiles"
	"github.com/hypermodeinc/modus/runtime/graphql"
	"github.com/hypermodeinc/modus/runtime/hostfunctions"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/middleware"
	"github.com/hypermodeinc/modus/runtime/neo4jclient"
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/hypermodeinc/modus/runtime/sqlclient"
	"github.com/hypermodeinc/modus/runtime/storage"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
)

// Starts any services that need to be started when the runtime starts.
func Start(ctx context.Context) context.Context {

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
	neo4jclient.Initialize()
	secrets.Initialize(ctx)
	storage.Initialize(ctx)
	db.Initialize(ctx)
	db.InitModusDb(ctx)
	actors.Initialize(ctx)
	graphql.Initialize()
	envfiles.MonitorEnvFiles(ctx)
	manifestdata.MonitorManifestFile(ctx)
	pluginmanager.Initialize(ctx)

	return ctx
}

// Stops any services that need to be stopped when the runtime stops.
func Stop(ctx context.Context) {

	// Actors need to be stopped first, as they may need to call wasm functions during shutdown.
	actors.Shutdown(ctx)

	// Now stop the wasm host.
	wasmhost.GetWasmHost(ctx).Close(ctx)

	// Stop the rest of the background services.
	// NOTE: Stopping services also has an order dependency.
	// If you need to change the order or add new services, be sure to test thoroughly.
	// Unlike start, these should each block until they are fully stopped.

	middleware.Shutdown()
	sqlclient.Shutdown()
	dgraphclient.ShutdownConns()
	neo4jclient.CloseDrivers(ctx)
	db.Stop(ctx)
	db.CloseModusDb(ctx)

	logger.Info(ctx).Msg("Shutdown complete.")
}
