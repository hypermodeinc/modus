/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"context"
	"fmt"
	"strings"

	"hypruntime/hostfunctions/compatibility"
	"hypruntime/logger"
	"hypruntime/plugins"
	"hypruntime/utils"
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
	for fnName, fnDef := range fnExports {
		if fnMeta, ok := plugin.Metadata.FnExports[fnName]; ok {

			plan, err := plugin.Planner.GetPlan(ctx, fnMeta, fnDef)
			if err != nil {
				logger.Err(ctx, err).
					Str("function", fnName).
					Str("plugin", plugin.Name()).
					Str("build_id", plugin.BuildId()).
					Msg("Error creating execution plan.")
				continue
			}

			fr.functions[fnName] = NewFunctionInfo(plugin, plan)
			names = append(names, fnName)

			logger.Info(ctx).
				Str("function", fnName).
				Str("plugin", plugin.Name()).
				Str("build_id", plugin.BuildId()).
				Msg("Registered exported function.")
		}
	}
	return names
}

func (fr *functionRegistry) RegisterImports(ctx context.Context, plugin *plugins.Plugin) []string {
	gaveCompatibilityWarning := false
	fnImports := plugin.Module.ImportedFunctions()
	names := make([]string, 0, len(fnImports))
	for _, fnDef := range fnImports {
		modName, fnName, _ := fnDef.Import()
		impName := fmt.Sprintf("%s.%s", modName, fnName)
		if _, ok := fr.functions[impName]; ok {
			names = append(names, impName)
			continue // already registered
		}

		fnMeta, found := plugin.Metadata.FnImports[impName]
		if !found {
			if modName == "hypermode" {
				if !gaveCompatibilityWarning {
					logger.Warn(ctx).
						Str("plugin", plugin.Name()).
						Str("build_id", plugin.BuildId()).
						Msg("Hypermode function imports are missing from the metadata. Using compatibility shims. Please update your SDK to the latest version.")
					gaveCompatibilityWarning = true
				}
				if m, err := compatibility.GetImportMetadataShim(impName); err == nil {
					fnMeta = m
				} else {
					logger.Err(ctx, err).
						Str("function", impName).
						Str("plugin", plugin.Name()).
						Str("build_id", plugin.BuildId()).
						Msg("Error creating compatibility shim.")
					continue
				}
			} else if fr.shouldIgnoreModule(modName) {
				continue
			} else {
				logger.Warn(ctx).
					Str("function", impName).
					Str("plugin", plugin.Name()).
					Str("build_id", plugin.BuildId()).
					Msg("Function is not registered in metadata imports.  The plugin may not work as expected.")
				continue
			}
		}

		plan, err := plugin.Planner.GetPlan(ctx, fnMeta, fnDef)
		if err != nil {
			logger.Err(ctx, err).
				Str("function", impName).
				Str("plugin", plugin.Name()).
				Str("build_id", plugin.BuildId()).
				Msg("Error creating execution plan.")
			continue
		}

		fr.functions[impName] = NewFunctionInfo(plugin, plan)
		names = append(names, impName)
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
			logger.Info(ctx).
				Str("function", name).
				Str("plugin", fnInfo.Plugin().Name()).
				Str("build_id", fnInfo.Plugin().BuildId()).
				Msg("Unregistered function.")
		}
	}
}

func (fr *functionRegistry) shouldIgnoreModule(name string) bool {
	switch name {
	case "wasi_snapshot_preview1", "wasi", "env", "runtime", "syscall", "test":
		return true
	}

	return strings.HasPrefix(name, "wasi_")
}
