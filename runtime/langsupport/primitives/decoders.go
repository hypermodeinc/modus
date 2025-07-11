/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package primitives

import (
	"math"
	"time"
)

func DecodeBool(val uint64) bool {
	return val != 0
}

func DecodeInt8(val uint64) int8 {
	return int8(val)
}

func DecodeInt16(val uint64) int16 {
	return int16(val)
}

func DecodeInt32(val uint64) int32 {
	return int32(val)
}

func DecodeInt64(val uint64) int64 {
	return int64(val)
}

func DecodeUint8(val uint64) uint8 {
	return uint8(val)
}

func DecodeUint16(val uint64) uint16 {
	return uint16(val)
}

func DecodeUint32(val uint64) uint32 {
	return uint32(val)
}

func DecodeUint64(val uint64) uint64 {
	return val
}

func DecodeInt(val uint64) int {
	return int(int32(val))
}

func DecodeUint(val uint64) uint {
	return uint(uint32(val))
}

func DecodeUintptr(val uint64) uintptr {
	return uintptr(uint32(val))
}

func DecodeFloat32(val uint64) float32 {
	return math.Float32frombits(uint32(val))
}

func DecodeFloat64(val uint64) float64 {
	return math.Float64frombits(val)
}

func DecodeDuration(val uint64) time.Duration {
	return time.Duration(val)
}
