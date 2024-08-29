/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"hypruntime/wasmhost"
)

type HypDSConfig struct {
	WasmHost *wasmhost.WasmHost
	MapTypes []string
}
