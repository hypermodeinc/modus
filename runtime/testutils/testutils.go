/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package testutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/functions"
	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
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

	fnInfo, ok := functions.NewFunctionInfo(name, f.Plugin, false)
	if !ok {
		return nil, fmt.Errorf("no function registered named %s", name)
	}

	params, err := functions.CreateParametersMap(fnInfo.Metadata(), paramValues...)
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

func NewWasmTestFixture(wasmFilePath string, customTypes map[string]reflect.Type, registrations []func(wasmhost.WasmHost) error) *WasmTestFixture {
	logger.Initialize()

	content, err := os.ReadFile(wasmFilePath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	host := wasmhost.NewWasmHost(ctx, registrations...)

	cm, err := host.CompileModule(ctx, content)
	if err != nil {
		panic(err)
	}

	md, err := metadata.GetMetadataFromCompiledModule(cm)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, utils.MetadataContextKey, md)
	ctx = context.WithValue(ctx, utils.CustomTypesContextKey, customTypes)

	filename := filepath.Base(wasmFilePath)
	plugin, err := plugins.NewPlugin(ctx, cm, filename, md)
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)

	registry := host.GetFunctionRegistry()
	registry.RegisterImports(ctx, plugin)

	return &WasmTestFixture{
		Context:     ctx,
		WasmHost:    host,
		Plugin:      plugin,
		customTypes: customTypes,
	}
}
