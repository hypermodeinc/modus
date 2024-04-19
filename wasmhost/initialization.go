/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func InitWasmRuntime(ctx context.Context) error {
	cfg := wazero.NewRuntimeConfig()
	cfg = cfg.WithCloseOnContextDone(true)
	RuntimeInstance = wazero.NewRuntimeWithConfig(ctx, cfg)
	return instantiateWasiFunctions(ctx)
}

func instantiateWasiFunctions(ctx context.Context) error {
	b := RuntimeInstance.NewHostModuleBuilder(wasi_snapshot_preview1.ModuleName)
	wasi_snapshot_preview1.NewFunctionExporter().ExportFunctions(b)

	// If we ever need to override any of the WASI functions, we can do so here.

	_, err := b.Instantiate(ctx)
	if err != nil {
		return fmt.Errorf("failed to instantiate the %s module: %w", wasi_snapshot_preview1.ModuleName, err)
	}

	return nil
}
