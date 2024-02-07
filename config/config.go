/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import (
	"hmruntime/dgraph"

	"github.com/tetratelabs/wazero"
)

// command line flag variables
var DgraphUrl *string
var PluginsPath *string

// map that holds the compiled modules for each plugin
var CompiledModules = make(map[string]wazero.CompiledModule)

// map that holds the function info for each resolver
var FunctionsMap = make(map[string]dgraph.FunctionInfo)

// Channel used to signal that registration is needed
var Register chan bool = make(chan bool)

// channel and flag used to signal the HTTP server
var ServerReady chan bool = make(chan bool)
var ServerWaiting = true

var WasmRuntime wazero.Runtime
