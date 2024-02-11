/*
 * Copyright 2023 Hypermode, Inc.
 */

package main

import (
	"context"
	"fmt"
	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/functions"
	"hmruntime/host"
	"hmruntime/plugins"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	config.ParseCommandLineFlags()

	// Ensure the plugins directory exists.
	if _, err := os.Stat(config.PluginsPath); os.IsNotExist(err) {
		err := os.MkdirAll(config.PluginsPath, 0755)
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to create plugins directory: %w", err))
		}
	}

	// Initialize the AWS configuration
	err := aws.Initialize(ctx)
	if err != nil {
		log.Println(err)
		log.Println("AWS functionality will be disabled")
	}

	// Initialize the WebAssembly runtime
	host.WasmRuntime, err = plugins.InitWasmRuntime(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer host.WasmRuntime.Close(ctx)

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

	err = startServer()
	log.Fatalln(err)

	// TODO: Shutdown gracefully
}
