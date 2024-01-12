/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// channel and flag used to signal the HTTP server
var serverReady chan bool = make(chan bool)
var serverWaiting = true

var dgraphUrl *string
var pluginsPath *string

func main() {
	ctx := context.Background()

	// Parse command-line flags
	var port = flag.Int("port", 8686, "The HTTP port to listen on.")
	dgraphUrl = flag.String("dgraph", "http://localhost:8080", "The Dgraph url to connect to.")
	pluginsPath = flag.String("plugins", "./plugins", "The path to the plugins directory.")
	flag.Parse()

	// Initialize the WebAssembly runtime
	var err error
	wasmRuntime, err = initWasmRuntime(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer wasmRuntime.Close(ctx)

	// Load plugins
	err = loadPlugins(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Watch for registration requests
	monitorRegistration(ctx)

	// Watch for schema changes
	monitorGqlSchema(ctx)

	// Watch for plugin changes
	err = watchPluginDirectory(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Start the HTTP server when we're ready
	<-serverReady
	serverWaiting = false
	fmt.Printf("Listening on port %d...\n", *port)
	http.HandleFunc("/graphql-worker", handleRequest)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	log.Fatalln(err)

	// TODO: Shutdown gracefully
}
