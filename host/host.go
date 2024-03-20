/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import "github.com/tetratelabs/wazero"

// runtime instance for the WASM modules
var WasmRuntime wazero.Runtime

// Channel used to signal that registration is needed
var RegistrationRequest chan bool = make(chan bool)
