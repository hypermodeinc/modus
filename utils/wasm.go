/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import wasm "github.com/tetratelabs/wazero/api"

func CopyMemory(mem wasm.Memory, sourceOffset, targetOffset, size uint32) bool {
	data, ok := mem.Read(sourceOffset, size)
	if !ok {
		return false
	}
	return mem.Write(targetOffset, data)
}
