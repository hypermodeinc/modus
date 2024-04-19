/*
 * Copyright 2024 Hypermode, Inc.
 */

package testutils

import (
	"context"
	_ "embed"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed data/test-as/test-as.wasm
var testWasm []byte

type WasmTestFixture struct {
	Context context.Context
	Runtime wazero.Runtime
	Module  api.Module
	Memory  api.Memory
}

func (f *WasmTestFixture) Close() error {
	return f.Runtime.Close(f.Context)
}

func NewWasmTestFixture() WasmTestFixture {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// TODO: refactor to share config with the real code so we're testing the same thing
	cfg := wazero.NewModuleConfig().
		WithSysWalltime().WithSysNanotime()

	mod, err := r.InstantiateWithConfig(ctx, testWasm, cfg)
	if err != nil {
		panic(err)
	}

	mem := mod.Memory()

	return WasmTestFixture{ctx, r, mod, mem}
}
