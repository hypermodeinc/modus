/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

func GetAlignmentPadding(offset, size uint32) uint32 {

	// maximum alignment is 4 bytes on 32-bit wasm
	if size > 4 {
		size = 4
	}

	mask := size - 1
	if offset&mask != 0 {
		return (offset | mask) + 1 - offset
	}

	return 0
}
