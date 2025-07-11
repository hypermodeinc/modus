/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
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
