/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import (
	"hmruntime/plugins"

	"github.com/tetratelabs/wazero"
)

// runtime instance for the WASM modules
var WasmRuntime wazero.Runtime

// Channel used to signal that registration is needed
var RegistrationRequest chan bool = make(chan bool)

// Global, thread-safe registry of all plugins loaded by the host
var Plugins = plugins.NewPluginRegistry()
