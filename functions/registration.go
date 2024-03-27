/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"reflect"
	"strings"

	"hmruntime/host"
	"hmruntime/logger"
	"hmruntime/schema"
)

// map that holds the function info for each resolver
var FunctionsMap = make(map[string]schema.FunctionInfo)

func MonitorRegistration(ctx context.Context) {
	go func() {
		for {
			select {
			case <-host.RegistrationRequest:
				err := registerFunctions(ctx, gqlSchema)
				if err != nil {
					logger.Err(ctx, err).Msg("Failed to register functions.")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func registerFunctions(ctx context.Context, gqlSchema string) error {

	// Get the function schema from the GraphQL schema.
	funcSchemas, err := schema.GetFunctionSchema(gqlSchema)
	if err != nil {
		return err
	}

	// Build a map of resolvers to function info, including the plugin name.
	// If there are function name conflicts between plugins, the last plugin loaded wins.
	var plugins = host.Plugins.GetAll()
	for _, plugin := range plugins {
		for _, scma := range funcSchemas {
			module := *plugin.Module
			for _, fn := range module.ExportedFunctions() {
				fnName := fn.ExportNames()[0]
				if strings.EqualFold(fnName, scma.FunctionName()) {
					info := schema.FunctionInfo{Plugin: &plugin, Schema: scma}
					resolver := scma.Resolver()
					oldInfo, existed := FunctionsMap[resolver]
					if existed && reflect.DeepEqual(oldInfo, info) {
						continue
					}
					FunctionsMap[resolver] = info

					logger.Info(ctx).
						Str("resolver", resolver).
						Str("function", fnName).
						Str("plugin", plugin.Name()).
						Str("build_id", plugin.BuildId()).
						Msg("Registered function.")
				}
			}
		}
	}

	// Cleanup any previously registered functions that are no longer in the schema or loaded modules.
	for resolver, info := range FunctionsMap {
		foundSchema := false
		for _, schema := range funcSchemas {
			if strings.EqualFold(info.FunctionName(), schema.FunctionName()) {
				foundSchema = true
				break
			}
		}
		_, foundPlugin := host.Plugins.GetByName(info.Plugin.Name())
		if !foundSchema || !foundPlugin {
			delete(FunctionsMap, resolver)
			logger.Info(ctx).
				Str("resolver", resolver).
				Str("function", info.FunctionName()).
				Str("plugin", info.Plugin.Name()).
				Str("build_id", info.Plugin.BuildId()).
				Msg("Unregistered function.")
		}
	}

	return nil
}
