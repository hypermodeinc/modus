/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import "github.com/tetratelabs/wazero"

// runtime instance for the WASM modules
var WasmRuntime wazero.Runtime

// map that holds all of the loaded plugins, indexed by their name
var Plugins = make(map[string]Plugin)

// Channel used to signal that registration is needed
var RegistrationRequest chan bool = make(chan bool)
