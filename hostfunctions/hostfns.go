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

	wasm "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const hypermodeHostModuleName string = "hypermode"

var hostFunctions []*hostFunctionDefinition

func addHostFunction(name string, function wasm.GoModuleFunc, options ...func(*hostFunctionDefinition)) {
	f := &hostFunctionDefinition{name: name, function: function}
	for _, option := range options {
		option(f)
	}
	hostFunctions = append(hostFunctions, f)
}

func withParams(qty int, vt wasm.ValueType) func(*hostFunctionDefinition) {
	return func(f *hostFunctionDefinition) {
		a := make([]wasm.ValueType, qty)
		for i := 0; i < qty; i++ {
			a[i] = vt
		}
		f.params = append(f.params, a...)
	}
}

func withResults(qty int, vt wasm.ValueType) func(*hostFunctionDefinition) {
	return func(f *hostFunctionDefinition) {
		a := make([]wasm.ValueType, qty)
		for i := 0; i < qty; i++ {
			a[i] = vt
		}
		f.results = append(f.results, a...)
	}
}

//nolint:unused
func withI32Param() func(*hostFunctionDefinition) {
	return withParams(1, wasm.ValueTypeI32)
}

//nolint:unused
func withI32Params(qty int) func(*hostFunctionDefinition) {
	return withParams(qty, wasm.ValueTypeI32)
}

//nolint:unused
func withI64Param() func(*hostFunctionDefinition) {
	return withParams(1, wasm.ValueTypeI64)
}

//nolint:unused
func withI64Params(qty int) func(*hostFunctionDefinition) {
	return withParams(qty, wasm.ValueTypeI64)
}

//nolint:unused
func withF32Param() func(*hostFunctionDefinition) {
	return withParams(1, wasm.ValueTypeF32)
}

//nolint:unused
func withF32Params(qty int) func(*hostFunctionDefinition) {
	return withParams(qty, wasm.ValueTypeF32)
}

//nolint:unused
func withF64Param() func(*hostFunctionDefinition) {
	return withParams(1, wasm.ValueTypeF64)
}

//nolint:unused
func withF64Params(qty int) func(*hostFunctionDefinition) {
	return withParams(qty, wasm.ValueTypeF64)
}

//nolint:unused
func withI32Result() func(*hostFunctionDefinition) {
	return withResults(1, wasm.ValueTypeI32)
}

//nolint:unused
func withI32Results(qty int) func(*hostFunctionDefinition) {
	return withResults(qty, wasm.ValueTypeI32)
}

//nolint:unused
func withI64Result() func(*hostFunctionDefinition) {
	return withResults(1, wasm.ValueTypeI64)
}

//nolint:unused
func withI64Results(qty int) func(*hostFunctionDefinition) {
	return withResults(qty, wasm.ValueTypeI64)
}

//nolint:unused
func withF32Result() func(*hostFunctionDefinition) {
	return withResults(1, wasm.ValueTypeF32)
}

//nolint:unused
func withF32Results(qty int) func(*hostFunctionDefinition) {
	return withResults(qty, wasm.ValueTypeF32)
}

//nolint:unused
func withF64Result() func(*hostFunctionDefinition) {
	return withResults(1, wasm.ValueTypeF64)
}

//nolint:unused
func withF64Results(qty int) func(*hostFunctionDefinition) {
	return withResults(qty, wasm.ValueTypeF64)
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
