/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package primitives

import (
	"math"
	"time"
)

func EncodeBool(val bool) uint64 {
	if val {
		return 1
	} else {
		return 0
	}
}

func EncodeInt8(val int8) uint64 {
	return uint64(uint32(val))
}

func EncodeInt16(val int16) uint64 {
	return uint64(uint32(val))
}

func EncodeInt32(val int32) uint64 {
	return uint64(uint32(val))
}

func EncodeInt64(val int64) uint64 {
	return uint64(val)
}

func EncodeUint8(val uint8) uint64 {
	return uint64(val)
}

func EncodeUint16(val uint16) uint64 {
	return uint64(val)
}

func EncodeUint32(val uint32) uint64 {
	return uint64(val)
}

func EncodeUint64(val uint64) uint64 {
	return val
}

func EncodeInt(val int) uint64 {
	return uint64(uint32(val))
}

func EncodeUint(val uint) uint64 {
	return uint64(val)
}

func EncodeUintptr(val uintptr) uint64 {
	return uint64(uint32(val))
}

func EncodeFloat32(val float32) uint64 {
	return uint64(math.Float32bits(val))
}

func EncodeFloat64(val float64) uint64 {
	return math.Float64bits(val)
}

func EncodeDuration(val time.Duration) uint64 {
	return uint64(val)
}
