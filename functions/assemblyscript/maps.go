/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"

	"hmruntime/plugins"

	"github.com/tetratelabs/wazero/api"
)

// Reference: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/map.ts

type kvp struct {
	Key   any `json:"key"`
	Value any `json:"value"`
}

func readMap(ctx context.Context, mem api.Memory, def plugins.TypeDefinition, offset uint32) ([]kvp, error) {

	// buckets, ok := mem.ReadUint32Le(offset)
	// if !ok {
	// 	return nil, fmt.Errorf("failed to read map buckets pointer")
	// }

	// bucketsMask, ok := mem.ReadUint32Le(offset + 4)
	// if !ok {
	// 	return nil, fmt.Errorf("failed to read map buckets mask")
	// }

	entries, ok := mem.ReadUint32Le(offset + 8)
	if !ok {
		return nil, fmt.Errorf("failed to read map entries pointer")
	}

	entriesCapacity, ok := mem.ReadUint32Le(offset + 12)
	if !ok {
		return nil, fmt.Errorf("failed to read map entries capacity")
	}

	// entriesOffset, ok := mem.ReadUint32Le(offset + 16)
	// if !ok {
	// 	return nil, fmt.Errorf("failed to read map entries offset")
	// }

	entriesCount, ok := mem.ReadUint32Le(offset + 20)
	if !ok {
		return nil, fmt.Errorf("failed to read map entries count")
	}

	// Handle empty map
	if entriesCount == 0 {
		return []kvp{}, nil
	}

	// the length of array buffer is stored 4 bytes before the offset
	byteLength, ok := mem.ReadUint32Le(entries - 4)
	if !ok {
		return nil, fmt.Errorf("failed to read map entries buffer length")
	}

	entrySize := byteLength / entriesCapacity
	keyType, valueType := getMapSubtypeInfo(def.Path)

	var valueOffset uint32
	switch keyType.Path {
	case "u64", "i64", "f64":
		// 64-bit keys have 8-byte value offset
		valueOffset = 8
	default:
		// everything else has 4-byte value offset
		valueOffset = 4
	}

	result := make([]kvp, entriesCount)
	for i := uint32(0); i < entriesCount; i++ {
		p := entries + (i * entrySize)
		k, err := readField(ctx, mem, keyType, p)
		if err != nil {
			return nil, err
		}
		v, err := readField(ctx, mem, valueType, p+valueOffset)
		if err != nil {
			return nil, err
		}

		result[i] = kvp{Key: k, Value: v}
	}

	return result, nil
}
