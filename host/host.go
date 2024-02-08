/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import "github.com/tetratelabs/wazero"

// runtime instance for the WASM modules
var WasmRuntime wazero.Runtime

// map that holds the compiled modules for each plugin
var CompiledModules = make(map[string]wazero.CompiledModule)
