/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package primitives

import (
	"time"

	wasm "github.com/tetratelabs/wazero/api"
)

func WriteBool(mem wasm.Memory, offset uint32, val bool) bool {
	b := byte(0)
	if val {
		b = 1
	}
	return mem.WriteByte(offset, b)
}

func WriteInt8(mem wasm.Memory, offset uint32, val int8) bool {
	return mem.WriteByte(offset, byte(val))
}

func WriteInt16(mem wasm.Memory, offset uint32, val int16) bool {
	return mem.WriteUint16Le(offset, uint16(val))
}

func WriteInt32(mem wasm.Memory, offset uint32, val int32) bool {
	return mem.WriteUint32Le(offset, uint32(val))
}

func WriteInt64(mem wasm.Memory, offset uint32, val int64) bool {
	return mem.WriteUint64Le(offset, uint64(val))
}

func WriteUint8(mem wasm.Memory, offset uint32, val uint8) bool {
	return mem.WriteByte(offset, val)
}

func WriteUint16(mem wasm.Memory, offset uint32, val uint16) bool {
	return mem.WriteUint16Le(offset, val)
}

func WriteUint32(mem wasm.Memory, offset uint32, val uint32) bool {
	return mem.WriteUint32Le(offset, val)
}

func WriteUint64(mem wasm.Memory, offset uint32, val uint64) bool {
	return mem.WriteUint64Le(offset, val)
}

func WriteInt(mem wasm.Memory, offset uint32, val int) bool {
	return mem.WriteUint32Le(offset, uint32(val))
}

func WriteUint(mem wasm.Memory, offset uint32, val uint) bool {
	return mem.WriteUint32Le(offset, uint32(val))
}

func WriteUintptr(mem wasm.Memory, offset uint32, val uintptr) bool {
	return mem.WriteUint32Le(offset, uint32(val))
}

func WriteFloat32(mem wasm.Memory, offset uint32, val float32) bool {
	return mem.WriteFloat32Le(offset, val)
}

func WriteFloat64(mem wasm.Memory, offset uint32, val float64) bool {
	return mem.WriteFloat64Le(offset, val)
}

func WriteDuration(mem wasm.Memory, offset uint32, val time.Duration) bool {
	return mem.WriteUint64Le(offset, uint64(val))
}
