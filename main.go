/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"hmruntime/dgraph"
	"hmruntime/monitor"
	"hmruntime/plugins"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()

	// Parse command-line flags
	var port = flag.Int("port", 8686, "The HTTP port to listen on.")
	dgraph.DgraphUrl = flag.String("dgraph", "http://localhost:8080", "The Dgraph url to connect to.")

	plugins.PluginsPath = flag.String("plugins", "./plugins", "The path to the plugins directory.")
	flag.StringVar(plugins.PluginsPath, "plugin", "./plugins", "alias for -plugins")

	flag.Parse()

	// Ensure the plugins directory exists.
	if _, err := os.Stat(*plugins.PluginsPath); os.IsNotExist(err) {
		err := os.MkdirAll(*plugins.PluginsPath, 0755)
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to create plugins directory: %w", err))
		}
	}

	// Initialize the WebAssembly runtime
	var err error
	plugins.WasmRuntime, err = plugins.InitWasmRuntime(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer plugins.WasmRuntime.Close(ctx)

	// Load plugins
	err = plugins.LoadPlugins(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Watch for registration requests
	monitor.MonitorRegistration(ctx)

	// Watch for schema changes
	monitor.MonitorGqlSchema(ctx)

	// Watch for plugin changes
	err = plugins.WatchPluginDirectory(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Start the HTTP server when we're ready
	<-monitor.ServerReady
	monitor.ServerWaiting = false
	fmt.Printf("Listening on port %d...\n", *port)
	http.HandleFunc("/graphql-worker", handleRequest)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	log.Fatalln(err)

	// TODO: Shutdown gracefully
}
