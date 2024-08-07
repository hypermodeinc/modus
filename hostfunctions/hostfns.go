/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/sqlclient"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const hypermodeHostModuleName string = "hypermode"

var hostFunctions []*hostFunctionDefinition

func addHostFunction(f *hostFunctionDefinition) {
	hostFunctions = append(hostFunctions, f)
}

func RegisterHostFunctions(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	instantiateHypermodeHostFunctions(ctx)
	instantiateWasiHostFunctions(ctx)

	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		sqlclient.ShutdownPGPools()
		return nil
	})
}

func instantiateHypermodeHostFunctions(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	b := wasmhost.RuntimeInstance.NewHostModuleBuilder(hypermodeHostModuleName)
	for _, hf := range hostFunctions {
		b.NewFunctionBuilder().
			WithGoModuleFunction(hf.function, hf.params, hf.results).
			Export(hf.name)
	}

	if _, err := b.Instantiate(ctx); err != nil {
		logger.Fatal(ctx).Err(err).
			Str("module", hypermodeHostModuleName).
			Msg("Failed to instantiate the host module.  Exiting.")
	}
}

func instantiateWasiHostFunctions(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	b := wasmhost.RuntimeInstance.NewHostModuleBuilder(wasi.ModuleName)
	wasi.NewFunctionExporter().ExportFunctions(b)

	// If we ever need to override any of the WASI functions, we can do so here.

	if _, err := b.Instantiate(ctx); err != nil {
		logger.Fatal(ctx).Err(err).
			Str("module", wasi.ModuleName).
			Msg("Failed to instantiate the host module.  Exiting.")
	}
}
