/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"fmt"
	"hmruntime/host"
	"hmruntime/schema"
	"log"
	"reflect"
	"strings"
)

// map that holds the function info for each resolver
var FunctionsMap = make(map[string]schema.FunctionInfo)

// channel used to signal when registration is completed
var RegistrationCompleted chan bool = make(chan bool)

func MonitorRegistration(ctx context.Context) {
	go func() {
		for {
			select {
			case <-host.RegistrationRequest:
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
	funcSchemas, err := schema.GetFunctionSchema(gqlSchema)
	if err != nil {
		return err
	}

	// Build a map of resolvers to function info, including the plugin name.
	// If there are function name conflicts between plugins, the last plugin loaded wins.
	for pluginName, cm := range host.CompiledModules {
		for _, scma := range funcSchemas {
			for _, fn := range cm.ExportedFunctions() {
				fnName := fn.ExportNames()[0]
				if strings.EqualFold(fnName, scma.FunctionName()) {
					info := schema.FunctionInfo{PluginName: pluginName, Schema: scma}
					resolver := scma.Resolver()
					oldInfo, existed := FunctionsMap[resolver]
					if existed && reflect.DeepEqual(oldInfo, info) {
						continue
					}
					FunctionsMap[resolver] = info
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
	for resolver, info := range FunctionsMap {
		foundSchema := false
		for _, schema := range funcSchemas {
			if strings.EqualFold(info.FunctionName(), schema.FunctionName()) {
				foundSchema = true
				break
			}
		}
		_, foundModule := host.CompiledModules[info.PluginName]
		if !foundSchema || !foundModule {
			delete(FunctionsMap, resolver)
			fmt.Printf("Unregistered old function '%s' for resolver '%s'\n", info.FunctionName(), resolver)
		}
	}

	// Signal that registration is complete (non-blocking send to avoid deadlock if no one is waiting)
	select {
	case RegistrationCompleted <- true:
	default:
	}

	return nil
}
