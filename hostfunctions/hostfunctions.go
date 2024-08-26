/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hmruntime/dgraphclient"
	"hmruntime/languages"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/plugins"
	"hmruntime/sqlclient"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	wasm "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const hypermodeHostModuleName string = "hypermode"

var hostFunctions []*hostFunctionDefinition

type hostFunctionDefinition struct {
	name     string
	function wasm.GoModuleFunction
	params   []wasm.ValueType
	results  []wasm.ValueType
}

func RegisterHostFunctions(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	instantiateHypermodeHostFunctions(ctx)
	instantiateWasiHostFunctions(ctx)

	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		sqlclient.ShutdownPGPools()
		dgraphclient.ShutdownConns()
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

func addHostFunction(name string, function wasm.GoModuleFunc, options ...func(*hostFunctionDefinition)) {
	f := &hostFunctionDefinition{name: name, function: function}
	for _, option := range options {
		option(f)
	}
	hostFunctions = append(hostFunctions, f)
}

func readParams(ctx context.Context, mod wasm.Module, stack []uint64, params ...any) error {
	if len(params) != len(stack) {
		return fmt.Errorf("expected a stack of size %d, but got %d", len(params), len(stack))
	}

	adapter, err := getWasmAdapter(ctx)
	if err != nil {
		return err
	}

	errs := make([]error, 0, len(params))
	for i, p := range params {
		if err := adapter.DecodeValue(ctx, mod, stack[i], p); err != nil {
			errs = append(errs, err)
		}
	}

	// clear the stack so that no input values can be treated as a result
	for i := 0; i < len(stack); i++ {
		stack[i] = 0
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func writeResults(ctx context.Context, mod wasm.Module, stack []uint64, results ...any) error {
	if len(results) > len(stack) {
		return fmt.Errorf("not enough stack space to write %d results", len(results))
	}

	adapter, err := getWasmAdapter(ctx)
	if err != nil {
		return err
	}

	errs := make([]error, 0, len(results))
	for i, r := range results {
		val, err := adapter.EncodeValue(ctx, mod, r)
		if err != nil {
			stack[i] = 0
			errs = append(errs, err)
		} else {
			stack[i] = val
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil

}

func getWasmAdapter(ctx context.Context) (languages.WasmAdapter, error) {
	p := plugins.GetPlugin(ctx)
	if p == nil {
		return nil, errors.New("no plugin found in context")
	}

	wa := p.Language.WasmAdapter()
	if wa == nil {
		return nil, errors.New("no wasm adapter found in plugin")
	}

	return wa, nil
}

// Each message is optional, but if provided, it will be logged at the appropriate time.
type hostFunctionMessages struct {
	Starting  string
	Completed string
	Cancelled string
	Error     string
	Detail    string
}

func callHostFunction(ctx context.Context, fn func() error, msgs *hostFunctionMessages) bool {

	if msgs != nil && msgs.Starting != "" {
		l := logger.Info(ctx).Bool("user_visible", true)
		if msgs.Detail != "" {
			l.Str("detail", msgs.Detail)
		}
		l.Msg(msgs.Starting)
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	if errors.Is(err, context.Canceled) {
		if msgs != nil && msgs.Cancelled != "" {
			l := logger.Warn(ctx).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.Detail != "" {
				l.Str("detail", msgs.Detail)
			}
			l.Msg(msgs.Cancelled)
		}
		return false
	} else if err != nil {
		if msgs != nil && msgs.Error != "" {
			l := logger.Err(ctx, err).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.Detail != "" {
				l.Str("detail", msgs.Detail)
			}
			l.Msg(msgs.Error)
		}
		return false
	} else {
		if msgs != nil && msgs.Completed != "" {
			l := logger.Info(ctx).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.Detail != "" {
				l.Str("detail", msgs.Detail)
			}
			l.Msg(msgs.Completed)
		}
		return true
	}
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
