/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"github.com/hypermodeinc/modus/runtime/wasmhost"
)

type HypDSConfig struct {
	WasmHost wasmhost.WasmHost
	MapTypes []string
}
