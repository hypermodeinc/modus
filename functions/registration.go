/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"

	"hmruntime/host"
	"hmruntime/logger"
	"hmruntime/plugins"
)

var Functions = make(map[string]FunctionInfo)
var TypeDefinitions = make(map[string]TypeInfo)

type FunctionInfo struct {
	Function plugins.FunctionSignature
	Plugin   *plugins.Plugin
}

type TypeInfo struct {
	Type   plugins.TypeDefinition
	Plugin *plugins.Plugin
}

func MonitorRegistration(ctx context.Context) {
	go func() {
		for {
			select {
			case <-host.RegistrationRequest:
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
	var plugins = host.Plugins.GetAll()
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

	// Save types from the metadata to the type definitions map
	for _, td := range plugin.Metadata.Types {
		TypeDefinitions[td.Name] = TypeInfo{
			Type:   td,
			Plugin: plugin,
		}
		r.types[td.Name] = true

		logger.Info(ctx).
			Str("type", td.Name).
			Str("plugin", plugin.Name()).
			Str("build_id", plugin.BuildId()).
			Msg("Registered custom data type.")
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

	// Cleanup any previously registered types
	for name, fn := range TypeDefinitions {
		if !r.types[name] {
			delete(TypeDefinitions, name)
			logger.Info(ctx).
				Str("type", name).
				Str("plugin", fn.Plugin.Name()).
				Str("build_id", fn.Plugin.BuildId()).
				Msg("Unregistered custom data type.")
		}
	}
}
