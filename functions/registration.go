/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"hmruntime/host"
	"hmruntime/schema"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

// map that holds the function info for each resolver
var FunctionsMap = make(map[string]schema.FunctionInfo)

// channel used to signal when registration is completed
var RegistrationCompleted chan bool = make(chan bool, 1)

func MonitorRegistration(ctx context.Context) {
	go func() {
		for {
			select {
			case <-host.RegistrationRequest:
				log.Info().Msg("Registering functions.")
				err := registerFunctions(gqlSchema)
				if err != nil {
					log.Err(err).Msg("Failed to register functions.")
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

					log.Info().
						Str("resolver", resolver).
						Str("function", fnName).
						Str("plugin", pluginName).
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
		_, foundModule := host.CompiledModules[info.PluginName]
		if !foundSchema || !foundModule {
			delete(FunctionsMap, resolver)
			log.Info().
				Str("resolver", resolver).
				Str("function", info.FunctionName()).
				Str("plugin", info.PluginName).
				Msg("Unregistered function.")
		}
	}

	// Signal that registration is complete.  This is a non-blocking send to
	// avoid a deadlock if no one is waiting, but the channel has a buffer size
	// of 1, so it will not lose the message if the receiver is slow to start.
	select {
	case RegistrationCompleted <- true:
	default:
	}

	return nil
}
