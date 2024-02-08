/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import (
	"github.com/tetratelabs/wazero"
)

// command line flag variables
var DgraphUrl string
var PluginsPath string
var NoReload bool

// map that holds the compiled modules for each plugin
var CompiledModules = make(map[string]wazero.CompiledModule)

// Channel used to signal that registration is needed
var Register chan bool = make(chan bool)

// channel and flag used to signal the HTTP server
var ServerReady chan bool = make(chan bool)
var ServerWaiting = true

var WasmRuntime wazero.Runtime
