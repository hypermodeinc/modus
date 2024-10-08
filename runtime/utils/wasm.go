/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
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
