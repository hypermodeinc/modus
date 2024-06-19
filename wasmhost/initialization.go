/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"fmt"

	"hmruntime/hostfunctions"
	"hmruntime/logger"
	"hmruntime/modules"
	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func InitWasmHost(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	cfg := wazero.NewRuntimeConfig()
	cfg = cfg.WithCloseOnContextDone(true)
	modules.RuntimeInstance = wazero.NewRuntimeWithConfig(ctx, cfg)

	err := initHostModules(ctx)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to initialize the WebAssembly runtime.  Exiting.")
	}
}

func initHostModules(ctx context.Context) error {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	// Instantiate the WASI functions
	err := instantiateWasiFunctions(ctx)
	if err != nil {
		return err
	}

	// Connect Hypermode host functions
	err = hostfunctions.Instantiate(ctx, &modules.RuntimeInstance)
	if err != nil {
		return err
	}

	return nil
}

func instantiateWasiFunctions(ctx context.Context) error {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	b := modules.RuntimeInstance.NewHostModuleBuilder(wasi.ModuleName)
	wasi.NewFunctionExporter().ExportFunctions(b)

	// If we ever need to override any of the WASI functions, we can do so here.

	_, err := b.Instantiate(ctx)
	if err != nil {
		return fmt.Errorf("failed to instantiate the %s module: %w", wasi.ModuleName, err)
	}

	return nil
}
