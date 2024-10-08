/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package primitives

import (
	"time"

	wasm "github.com/tetratelabs/wazero/api"
)

func ReadBool(mem wasm.Memory, offset uint32) (bool, bool) {
	val, ok := mem.ReadByte(offset)
	return val != 0, ok
}

func ReadInt8(mem wasm.Memory, offset uint32) (int8, bool) {
	val, ok := mem.ReadByte(offset)
	return int8(val), ok
}

func ReadInt16(mem wasm.Memory, offset uint32) (int16, bool) {
	val, ok := mem.ReadUint16Le(offset)
	return int16(val), ok
}

func ReadInt32(mem wasm.Memory, offset uint32) (int32, bool) {
	val, ok := mem.ReadUint32Le(offset)
	return int32(val), ok
}

func ReadInt64(mem wasm.Memory, offset uint32) (int64, bool) {
	val, ok := mem.ReadUint64Le(offset)
	return int64(val), ok
}

func ReadUint8(mem wasm.Memory, offset uint32) (uint8, bool) {
	val, ok := mem.ReadByte(offset)
	return val, ok
}

func ReadUint16(mem wasm.Memory, offset uint32) (uint16, bool) {
	val, ok := mem.ReadUint16Le(offset)
	return val, ok
}

func ReadUint32(mem wasm.Memory, offset uint32) (uint32, bool) {
	val, ok := mem.ReadUint32Le(offset)
	return val, ok
}

func ReadUint64(mem wasm.Memory, offset uint32) (uint64, bool) {
	val, ok := mem.ReadUint64Le(offset)
	return val, ok
}

func ReadInt(mem wasm.Memory, offset uint32) (int, bool) {
	val, ok := mem.ReadUint32Le(offset)
	return int(int32(val)), ok
}

func ReadUint(mem wasm.Memory, offset uint32) (uint, bool) {
	val, ok := mem.ReadUint32Le(offset)
	return uint(uint32(val)), ok
}

func ReadUintptr(mem wasm.Memory, offset uint32) (uintptr, bool) {
	val, ok := mem.ReadUint32Le(offset)
	return uintptr(uint32(val)), ok
}

func ReadFloat32(mem wasm.Memory, offset uint32) (float32, bool) {
	val, ok := mem.ReadFloat32Le(offset)
	return val, ok
}

func ReadFloat64(mem wasm.Memory, offset uint32) (float64, bool) {
	val, ok := mem.ReadFloat64Le(offset)
	return val, ok
}

func ReadDuration(mem wasm.Memory, offset uint32) (time.Duration, bool) {
	val, ok := mem.ReadUint64Le(offset)
	return time.Duration(val), ok
}
