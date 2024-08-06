/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"fmt"

	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/plugins/metadata"
)

var functions = make(map[string]*functionInfo)

type functionInfo struct {
	Function *metadata.Function
	Plugin   *plugins.Plugin
}

func GetFunction(fnName string) (*metadata.Function, error) {
	f, _, err := GetFunctionAndPlugin(fnName)
	return f, err
}

func GetFunctionAndPlugin(fnName string) (*metadata.Function, *plugins.Plugin, error) {
	info, ok := functions[fnName]
	if !ok {
		return nil, nil, fmt.Errorf("no function registered named %s", fnName)
	}
	return info.Function, info.Plugin, nil
}

func RegisterFunctions(ctx context.Context, plugins []*plugins.Plugin) {
	r := &registration{
		functions: make(map[string]bool),
		types:     make(map[string]bool),
	}

	for _, plugin := range plugins {
		r.registerPlugin(ctx, plugin)
	}

	r.cleanup(ctx)

	triggerFunctionsLoaded(ctx)
}

type registration struct {
	functions map[string]bool
	types     map[string]bool
}

func (r *registration) registerPlugin(ctx context.Context, plugin *plugins.Plugin) {

	// Save exported functions from the metadata to the functions map
	for _, fn := range plugin.Metadata.FnExports {
		functions[fn.Name] = &functionInfo{
			Function: fn,
			Plugin:   plugin,
		}
		r.functions[fn.Name] = true

		logger.Info(ctx).
			Str("function", fn.Name).
			Str("plugin", plugin.Name()).
			Str("build_id", plugin.BuildId()).
			Msg("Registered function.")
	}
}

func (r *registration) cleanup(ctx context.Context) {

	// Cleanup any previously registered functions
	for name, fn := range functions {
		if !r.functions[name] {
			delete(functions, name)
			logger.Info(ctx).
				Str("function", name).
				Str("plugin", fn.Plugin.Name()).
				Str("build_id", fn.Plugin.BuildId()).
				Msg("Unregistered function.")
		}
	}
}
