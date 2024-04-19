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

func readArray(ctx context.Context, mem api.Memory, def plugins.TypeDefinition, offset uint32) ([]any, error) {

	dataStart, ok := mem.ReadUint32Le(offset + 4)
	if !ok {
		return nil, fmt.Errorf("failed to read array data start pointer")
	}

	// byteLength, ok := mem.ReadUint32Le(offset + 8)
	// if !ok {
	// 	return nil, fmt.Errorf("failed to read array bytes length")
	// }

	arrLen, ok := mem.ReadUint32Le(offset + 12)
	if !ok {
		return nil, fmt.Errorf("failed to read array length")
	}

	// Handle empty array
	if arrLen == 0 {
		return []any{}, nil
	}

	itemType := getArraySubtypeInfo(def.Path)

	var itemSize uint32
	switch itemType.Path {
	case "u64", "i64", "f64":
		itemSize = 8
	case "u32", "i32", "f32":
		itemSize = 4
	case "u16", "i16":
		itemSize = 2
	case "u8", "i8", "bool":
		itemSize = 1
	default:
		itemSize = 4 // pointer
	}

	result := make([]any, arrLen)
	for i := uint32(0); i < arrLen; i++ {
		o, err := readField(ctx, mem, itemType, dataStart+(i*itemSize))
		if err != nil {
			return nil, err
		}
		result[i] = o
	}

	return result, nil
}
