/*
 * Copyright 2024 Hypermode, Inc.
 */

package testutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"hypruntime/functions"
	"hypruntime/plugins"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
	"hypruntime/wasmhost"

	wasm "github.com/tetratelabs/wazero/api"
)

type WasmTestFixture struct {
	Context        context.Context
	WasmHost       *wasmhost.WasmHost
	Plugin         *plugins.Plugin
	Buffers        *utils.OutputBuffers
	Module         wasm.Module
	customTypes    map[string]reflect.Type
	customTypesRev map[reflect.Type]string
}

func (f *WasmTestFixture) Close() {
	f.WasmHost.Close(f.Context)
}

func (f *WasmTestFixture) AddCustomType(name string, typ reflect.Type) {
	if f.customTypes == nil {
		f.customTypes = make(map[string]reflect.Type)
		f.customTypesRev = make(map[reflect.Type]string)
	}

	f.customTypes[name] = typ
	f.customTypesRev[typ] = name
}

func (f *WasmTestFixture) InvokeFunction(name string, paramValues ...any) (any, error) {
	fn, ok := f.Plugin.Metadata.FnExports[name]
	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	params, err := functions.CreateParametersMap(fn, paramValues...)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(f.Context, utils.CustomTypesContextKey, f.customTypes)
	ctx = context.WithValue(ctx, utils.CustomTypesRevContextKey, f.customTypesRev)

	info, err := f.WasmHost.InvokeFunction(ctx, f.Plugin, fn, params)
	if err != nil {
		return nil, err
	}

	return info.Result, nil
}

type testContextKey struct{}

func GetTestT(ctx context.Context) *testing.T {
	return ctx.Value(testContextKey{}).(*testing.T)
}

func NewWasmTestFixture(wasmFilePath string, t *testing.T, hostOpts ...func(*wasmhost.WasmHost) error) *WasmTestFixture {
	content, err := os.ReadFile(wasmFilePath)
	if err != nil {
		panic(err)
	}

	ctx := context.WithValue(context.Background(), testContextKey{}, t)
	host := wasmhost.NewWasmHost(ctx, hostOpts...)

	cm, err := host.CompileModule(ctx, content)
	if err != nil {
		panic(err)
	}

	md, err := metadata.GetMetadata(ctx, cm)
	if err != nil {
		panic(err)
	}

	filename := filepath.Base(wasmFilePath)
	plugin := plugins.NewPlugin(cm, filename, md)
	buffers := &utils.OutputBuffers{}

	mod, err := host.GetModuleInstance(ctx, plugin, buffers)
	if err != nil {
		panic(err)
	}

	return &WasmTestFixture{
		Context:  ctx,
		WasmHost: host,
		Plugin:   plugin,
		Buffers:  buffers,
		Module:   mod,
	}
}
