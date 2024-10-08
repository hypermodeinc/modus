/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hash

import (
	"bytes"
	"encoding/binary"
	"math/bits"

	"github.com/OneOfOne/xxhash"
)

// We'll leverage the xxhash library where possible.

func GetHashCode(data any) uint32 {
	switch t := data.(type) {
	case []byte:
		return xxhash.Checksum32(t)
	case uint64, int64:
		buf := bytes.NewBuffer(make([]byte, 0, 8))
		_ = binary.Write(buf, binary.LittleEndian, data)
		return xxhash.Checksum32(buf.Bytes())
	case uint32, int32, uint, int:
		buf := bytes.NewBuffer(make([]byte, 0, 4))
		_ = binary.Write(buf, binary.LittleEndian, data)
		return xxhash.Checksum32(buf.Bytes())
	case uint16:
		return hash32(uint32(t), 2)
	case int16:
		return hash32(uint32(t), 2)
	case uint8:
		return hash32(uint32(t), 1)
	case int8:
		return hash32(uint32(t), 1)
	case bool:
		if t {
			return hash32(1, 1)
		} else {
			return hash32(0, 1)
		}
	}

	return 0
}

// Remaining code ported from https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/util/hash.ts
// Note that AssemblyScript's implementation differs from xxhash for values less than 4 bytes in length

// primes
// const xxH32_P1 uint32 = 2654435761
const xxH32_P2 uint32 = 2246822519
const xxH32_P3 uint32 = 3266489917
const xxH32_P4 uint32 = 668265263
const xxH32_P5 uint32 = 374761393
const xxH32_SEED uint32 = 0

func hash32(val uint32, size uint32) uint32 {
	h := xxH32_SEED + xxH32_P5 + size
	h += val * xxH32_P3
	h = bits.RotateLeft32(h, 17) * xxH32_P4
	h ^= h >> 15
	h *= xxH32_P2
	h ^= h >> 13
	h *= xxH32_P3
	h ^= h >> 16
	return h
}
