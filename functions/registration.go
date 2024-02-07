/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"fmt"
	"hmruntime/config"
	"hmruntime/dgraph"
	"log"
	"reflect"
	"strings"
)

func MonitorRegistration(ctx context.Context) {
	go func() {
		for {
			select {
			case <-config.Register:
				log.Printf("Registering functions")
				err := registerFunctions(gqlSchema)
				if err != nil {
					log.Printf("Failed to register functions: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func registerFunctions(gqlSchema string) error {

	// Get the function schema from the GraphQL schema.
	funcSchemas, err := dgraph.GetFunctionSchema(gqlSchema)
	if err != nil {
		return err
	}

	// Build a map of resolvers to function info, including the plugin name.
	// If there are function name conflicts between plugins, the last plugin loaded wins.
	for pluginName, cm := range config.CompiledModules {
		for _, schema := range funcSchemas {
			for _, fn := range cm.ExportedFunctions() {
				fnName := fn.ExportNames()[0]
				if strings.EqualFold(fnName, schema.FunctionName()) {
					info := dgraph.FunctionInfo{PluginName: pluginName, Schema: schema}
					resolver := schema.Resolver()
					oldInfo, existed := config.FunctionsMap[resolver]
					if existed && reflect.DeepEqual(oldInfo, info) {
						continue
					}
					config.FunctionsMap[resolver] = info
					if existed {
						fmt.Printf("Re-registered %s to use %s in %s\n", resolver, fnName, pluginName)
					} else {
						fmt.Printf("Registered %s to use %s in %s\n", resolver, fnName, pluginName)
					}
				}
			}
		}
	}

	// Cleanup any previously registered functions that are no longer in the schema or loaded modules.
	for resolver, info := range config.FunctionsMap {
		foundSchema := false
		for _, schema := range funcSchemas {
			if strings.EqualFold(info.FunctionName(), schema.FunctionName()) {
				foundSchema = true
				break
			}
		}
		_, foundModule := config.CompiledModules[info.PluginName]
		if !foundSchema || !foundModule {
			delete(config.FunctionsMap, resolver)
			fmt.Printf("Unregistered old function '%s' for resolver '%s'\n", info.FunctionName(), resolver)
		}
	}

	// If the HTTP server is waiting, signal that we're ready.
	if config.ServerWaiting {
		config.ServerReady <- true
	}

	return nil
}
