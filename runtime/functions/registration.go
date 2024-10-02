/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package functions

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func NewFunctionRegistry() FunctionRegistry {
	return &functionRegistry{
		functions: make(map[string]FunctionInfo),
	}
}

type FunctionRegistry interface {
	GetFunctionInfo(fnName string) (FunctionInfo, error)
	RegisterAllFunctions(ctx context.Context, plugins ...*plugins.Plugin)
	RegisterImports(ctx context.Context, plugin *plugins.Plugin) []string
	RegisterExports(ctx context.Context, plugin *plugins.Plugin) []string
}

type functionRegistry struct {
	functions map[string]FunctionInfo
}

func (fr *functionRegistry) GetFunctionInfo(fnName string) (FunctionInfo, error) {
	info, ok := fr.functions[fnName]
	if !ok {
		return nil, fmt.Errorf("no function registered named %s", fnName)
	}
	return info, nil
}

func (fr *functionRegistry) RegisterAllFunctions(ctx context.Context, plugins ...*plugins.Plugin) {
	var names []string
	for _, plugin := range plugins {
		ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)
		ctx = context.WithValue(ctx, utils.MetadataContextKey, plugin.Metadata)
		importNames := fr.RegisterImports(ctx, plugin)
		exportNames := fr.RegisterExports(ctx, plugin)
		names = append(names, importNames...)
		names = append(names, exportNames...)
	}

	fr.cleanup(ctx, names)

	triggerFunctionsLoaded(ctx)
}

func (fr *functionRegistry) RegisterExports(ctx context.Context, plugin *plugins.Plugin) []string {
	fnExports := plugin.Module.ExportedFunctions()
	names := make([]string, 0, len(fnExports))
	for fnName := range fnExports {
		info, ok := NewFunctionInfo(fnName, plugin, false)
		if ok {
			fr.functions[fnName] = info
			names = append(names, fnName)

			logger.Info(ctx).
				Str("function", fnName).
				Str("plugin", plugin.Name()).
				Str("build_id", plugin.BuildId()).
				Msg("Registered function.")
		}
	}
	return names
}

func (fr *functionRegistry) RegisterImports(ctx context.Context, plugin *plugins.Plugin) []string {
	fnImports := plugin.Module.ImportedFunctions()
	names := make([]string, 0, len(fnImports))
	for _, fnDef := range fnImports {
		modName, fnName, _ := fnDef.Import()
		impName := fmt.Sprintf("%s.%s", modName, fnName)
		info, ok := NewFunctionInfo(impName, plugin, true)
		if ok {
			fr.functions[impName] = info
			names = append(names, impName)
		}
	}
	return names
}

func (fr *functionRegistry) cleanup(ctx context.Context, registeredNames []string) {
	m := make(map[string]bool, len(registeredNames))
	for _, name := range registeredNames {
		m[name] = true
	}

	for name, fnInfo := range fr.functions {
		if !m[name] {
			delete(fr.functions, name)

			if !fnInfo.IsImport() {
				logger.Info(ctx).
					Str("function", name).
					Str("plugin", fnInfo.Plugin().Name()).
					Str("build_id", fnInfo.Plugin().BuildId()).
					Msg("Unregistered function.")
			}
		}
	}
}
