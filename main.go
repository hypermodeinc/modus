/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"hmruntime/config"
	"hmruntime/dgraph"
	"hmruntime/functions"
	"hmruntime/plugins"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	// Parse command-line flags
	var port = flag.Int("port", 8686, "The HTTP port to listen on.")
	dgraph.DgraphUrl = flag.String("dgraph", "http://localhost:8080", "The Dgraph url to connect to.")

	config.PluginsPath = flag.String("plugins", "./plugins", "The path to the plugins directory.")
	flag.StringVar(config.PluginsPath, "plugin", "./plugins", "alias for -plugins")

	flag.Parse()

	// Ensure the plugins directory exists.
	if _, err := os.Stat(*config.PluginsPath); os.IsNotExist(err) {
		err := os.MkdirAll(*config.PluginsPath, 0755)
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to create plugins directory: %w", err))
		}
	}

	// Initialize the WebAssembly runtime
	var err error
	config.WasmRuntime, err = plugins.InitWasmRuntime(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer config.WasmRuntime.Close(ctx)

	// Load plugins
	err = plugins.LoadPlugins(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Watch for registration requests
	functions.MonitorRegistration(ctx)

	// Watch for schema changes
	functions.MonitorGqlSchema(ctx)

	// Watch for plugin changes
	err = plugins.WatchPluginDirectory(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	err = startServer(port)
	log.Fatalln(err)

	// TODO: Shutdown gracefully
}
