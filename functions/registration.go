/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/wasmhost"
)

var Functions = make(map[string]FunctionInfo)

type FunctionInfo struct {
	Function plugins.FunctionSignature
	Plugin   *plugins.Plugin
}

func MonitorRegistration(ctx context.Context) {
	go func() {
		for {
			select {
			case <-wasmhost.RegistrationRequest:
				r := newRegistration()
				r.registerAll(ctx)
				r.cleanup(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

type registration struct {
	functions map[string]bool
	types     map[string]bool
}

func newRegistration() *registration {
	return &registration{
		functions: make(map[string]bool),
		types:     make(map[string]bool),
	}
}

func (r *registration) registerAll(ctx context.Context) {
	var plugins = wasmhost.Plugins.GetAll()
	for _, plugin := range plugins {
		r.registerPlugin(ctx, &plugin)
	}
}

func (r *registration) registerPlugin(ctx context.Context, plugin *plugins.Plugin) {

	// Save functions from the metadata to the functions map
	for _, fn := range plugin.Metadata.Functions {
		Functions[fn.Name] = FunctionInfo{
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
	for name, fn := range Functions {
		if !r.functions[name] {
			delete(Functions, name)
			logger.Info(ctx).
				Str("function", name).
				Str("plugin", fn.Plugin.Name()).
				Str("build_id", fn.Plugin.BuildId()).
				Msg("Unregistered function.")
		}
	}
}
