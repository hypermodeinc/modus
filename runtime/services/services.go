/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package services

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/aws"
	"github.com/hypermodeinc/modus/runtime/collections"
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
	aws.Initialize(ctx)
	secrets.Initialize(ctx)
	storage.Initialize(ctx)
	db.Initialize(ctx)
	collections.Initialize(ctx)
	manifestdata.MonitorManifestFile(ctx)
	envfiles.MonitorEnvFiles(ctx)
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
	middleware.Shutdown()
	sqlclient.ShutdownPGPools()
	dgraphclient.ShutdownConns()
	neo4jclient.CloseDrivers(ctx)
	logger.Close()
	db.Stop(ctx)
}
