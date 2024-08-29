/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"

	"hypruntime/logger"
	"hypruntime/plugins"
	"hypruntime/utils"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmHost struct {
	runtime       wazero.Runtime
	hostFunctions []*HostFunction
}

func NewWasmHost(ctx context.Context, opts ...func(*WasmHost) error) *WasmHost {
	cfg := wazero.NewRuntimeConfig().WithCloseOnContextDone(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, cfg)
	wasi.MustInstantiate(ctx, runtime)
	host := &WasmHost{runtime: runtime}

	for _, reg := range opts {
		if err := reg(host); err != nil {
			logger.Fatal(ctx).Err(err).Msg("Failed to apply an option to the WASM host.")
			return nil
		}
	}

	if err := host.instantiateHostFunctions(ctx); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to instantiate host functions.")
		return nil
	}

	return host
}

func (host *WasmHost) Close(ctx context.Context) {
	if err := host.runtime.Close(ctx); err != nil {
		logger.Err(ctx, err).Msg("Failed to cleanly close the WASM runtime.")
	}
}

// Gets a module instance for the given plugin, used for a single invocation.
func (host *WasmHost) GetModuleInstance(ctx context.Context, plugin *plugins.Plugin, buffers *utils.OutputBuffers) (wasm.Module, error) {

	// Get the logger and writers for the plugin's stdout and stderr.
	log := logger.Get(ctx).With().Bool("user_visible", true).Logger()
	wInfoLog := logger.NewLogWriter(&log, zerolog.InfoLevel)
	wErrorLog := logger.NewLogWriter(&log, zerolog.ErrorLevel)

	// Capture stdout/stderr both to logs, and to provided writers.
	wOut := io.MultiWriter(&buffers.StdOut, wInfoLog)
	wErr := io.MultiWriter(&buffers.StdErr, wErrorLog)

	// Configure the module instance.
	// Note, we use an anonymous module name (empty string) here,
	// for concurrency and performance reasons.
	// See https://github.com/tetratelabs/wazero/pull/2275
	// And https://gophers.slack.com/archives/C040AKTNTE0/p1719587772724619?thread_ts=1719522663.531579&cid=C040AKTNTE0
	cfg := wazero.NewModuleConfig().
		WithName("").
		WithSysWalltime().WithSysNanotime().
		WithRandSource(rand.Reader).
		WithStdout(wOut).WithStderr(wErr)

	// Instantiate the plugin as a module.
	// NOTE: This will also invoke the plugin's `_start` function,
	// which will call any top-level code in the plugin.
	mod, err := host.runtime.InstantiateModule(ctx, plugin.Module, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the plugin module: %w", err)
	}

	return mod, nil
}

func (host *WasmHost) CompileModule(ctx context.Context, bytes []byte) (wazero.CompiledModule, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	cm, err := host.runtime.CompileModule(ctx, bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile the plugin: %w", err)
	}

	return cm, nil
}
