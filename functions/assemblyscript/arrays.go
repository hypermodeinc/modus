/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"hmruntime/plugins"

	wasm "github.com/tetratelabs/wazero/api"
)

func readArray(ctx context.Context, mem wasm.Memory, def plugins.TypeDefinition, offset uint32) ([]any, error) {

	dataStart, ok := mem.ReadUint32Le(offset + 4)
	if !ok {
		return nil, fmt.Errorf("failed to read array data start pointer")
	}

	byteLength, ok := mem.ReadUint32Le(offset + 8)
	if !ok {
		return nil, fmt.Errorf("failed to read array bytes length")
	}

	arrLen, ok := mem.ReadUint32Le(offset + 12)
	if !ok {
		return nil, fmt.Errorf("failed to read array length")
	}

	// Handle empty array to avoid division by zero
	if arrLen == 0 {
		return []any{}, nil
	}

	itemSize := byteLength / arrLen
	itemType := getArraySubtypeInfo(def.Path)

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
