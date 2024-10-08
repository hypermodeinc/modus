/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package index

import (
	"encoding/binary"
	"fmt"
	"math"
)

var (
	ErrVectorIndexAlreadyExists = fmt.Errorf("vector index already exists")
	ErrVectorIndexNotFound      = fmt.Errorf("vector index not found")
)

// BytesAsFloatArray(encoded) converts encoded into a []T,
// where T is either float32 or float64, depending on the value of floatBits.
// Let floatBytes = floatBits/8. If len(encoded) % floatBytes is
// not 0, it will ignore any trailing bytes, and simply convert floatBytes
// bytes at a time to generate the entries.
// The result is appended to the given retVal slice. If retVal is nil
// then a new slice is created and appended to.
func BytesAsFloatArray(encoded []byte, retVal *[]float32) {
	// Unfortunately, this is not as simple as casting the result,
	// and it is also not possible to directly use the
	// golang "unsafe" library to directly do the conversion.
	// The machine where this operation gets run might prefer
	// BigEndian/LittleEndian, but the machine that sent it may have
	// preferred the other, and there is no way to tell!
	//
	// The solution below, unfortunately, requires another memory
	// allocation.
	// TODO Potential optimization: If we detect that current machine is
	// using LittleEndian format, there might be a way of making this
	// work with the golang "unsafe" library.
	floatBytes := 4

	*retVal = (*retVal)[:0]
	resultLen := len(encoded) / floatBytes
	if resultLen == 0 {
		return
	}
	for i := 0; i < resultLen; i++ {
		// Assume LittleEndian for encoding since this is
		// the assumption elsewhere when reading from client.
		// See dgraph-io/dgo/protos/api.pb.go
		// See also dgraph-io/dgraph/types/conversion.go
		// This also seems to be the preference from many examples
		// I have found via Google search. It's unclear why this
		// should be a preference.
		if retVal == nil {
			retVal = &[]float32{}
		}
		*retVal = append(*retVal, BytesToFloat(encoded))

		encoded = encoded[(floatBytes):]
	}
}

func BytesToFloat(encoded []byte) float32 {
	bits := binary.LittleEndian.Uint32(encoded)
	return math.Float32frombits(bits)
}
