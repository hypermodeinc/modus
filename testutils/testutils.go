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
	"hypruntime/langsupport"
	"hypruntime/logger"
	"hypruntime/plugins"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
	"hypruntime/wasmhost"
)

type WasmTestFixture struct {
	Context     context.Context
	WasmHost    wasmhost.WasmHost
	Plugin      *plugins.Plugin
	customTypes map[string]reflect.Type
}

func (f *WasmTestFixture) Close() {
	f.WasmHost.Close(f.Context)
}

func (f *WasmTestFixture) AddCustomType(name string, typ reflect.Type) {
	f.customTypes[name] = typ
}

func (f *WasmTestFixture) CallFunction(t *testing.T, name string, paramValues ...any) (any, error) {
	ctx := context.WithValue(f.Context, testContextKey{}, t)

	fnMeta, ok := f.Plugin.Metadata.FnExports[name]
	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	fnDef := f.Plugin.Module.ExportedFunctions()[name]
	plan, err := f.Plugin.Planner.GetPlan(ctx, fnMeta, fnDef)
	if err != nil {
		return nil, err
	}

	fnInfo := functions.NewFunctionInfo(f.Plugin, plan)

	params, err := functions.CreateParametersMap(fnMeta, paramValues...)
	if err != nil {
		return nil, err
	}

	execInfo, err := f.WasmHost.CallFunction(ctx, fnInfo, params)
	if err != nil {
		return nil, err
	}

	return execInfo.Result(), nil
}

func (f *WasmTestFixture) NewPlanner() langsupport.Planner {
	return f.Plugin.Language.NewPlanner(f.Plugin.Metadata)
}

type testContextKey struct{}

func GetTestT(ctx context.Context) *testing.T {
	return ctx.Value(testContextKey{}).(*testing.T)
}

func NewWasmTestFixture(wasmFilePath string, hostOpts ...func(wasmhost.WasmHost) error) *WasmTestFixture {
	logger.Initialize()

	content, err := os.ReadFile(wasmFilePath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
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
	registry := host.GetFunctionRegistry()
	registry.RegisterImports(ctx, plugin)

	f := &WasmTestFixture{
		WasmHost:    host,
		Plugin:      plugin,
		customTypes: make(map[string]reflect.Type),
	}

	ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)
	ctx = context.WithValue(ctx, utils.MetadataContextKey, md)
	ctx = context.WithValue(ctx, utils.CustomTypesContextKey, f.customTypes)
	f.Context = ctx

	return f
}
