/*
 * Copyright 2024 Hypermode, Inc.
 */

package host

import (
	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Module *wazero.CompiledModule
}
