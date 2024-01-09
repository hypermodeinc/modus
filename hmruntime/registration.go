/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/tetratelabs/wazero"
)

// map that holds the compiled modules for each plugin
var compiledModules = make(map[string]wazero.CompiledModule)

// map that holds the function info for each resolver
var functionsMap = make(map[string]functionInfo)

// Channel used to signal that registration is needed
var register chan bool = make(chan bool)

func monitorRegistration(ctx context.Context) {
	go func() {
		for {
			select {
			case <-register:
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
	funcSchemas, err := getFunctionSchema(gqlSchema)
	if err != nil {
		return err
	}

	// Build a map of resolvers to function info, including the plugin name.
	// If there are function name conflicts between plugins, the last plugin loaded wins.
	for pluginName, cm := range compiledModules {
		for _, schema := range funcSchemas {
			for _, fn := range cm.ExportedFunctions() {
				fnName := fn.ExportNames()[0]
				if strings.EqualFold(fnName, schema.FunctionName()) {
					info := functionInfo{pluginName, schema}
					resolver := schema.Resolver()
					oldInfo, existed := functionsMap[resolver]
					if existed && reflect.DeepEqual(oldInfo, info) {
						continue
					}
					functionsMap[resolver] = info
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
	for resolver, info := range functionsMap {
		foundSchema := false
		for _, schema := range funcSchemas {
			if strings.EqualFold(info.FunctionName(), schema.FunctionName()) {
				foundSchema = true
				break
			}
		}
		_, foundModule := compiledModules[info.PluginName]
		if !foundSchema || !foundModule {
			delete(functionsMap, resolver)
			fmt.Printf("Unregistered old function '%s' for resolver '%s'\n", info.FunctionName(), resolver)
		}
	}

	// If the HTTP server is waiting, signal that we're ready.
	if serverWaiting {
		serverReady <- true
	}

	return nil
}
