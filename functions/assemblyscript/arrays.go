/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

// Reference: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/array.ts

func readArray(ctx context.Context, mem wasm.Memory, def plugins.TypeDefinition, offset uint32) ([]any, error) {

	// buffer, ok := mem.ReadUint32Le(offset)
	// if !ok {
	// 	return nil, fmt.Errorf("failed to read array buffer pointer")
	// }

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
	itemSize := getItemSize(itemType)

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

func writeArray(ctx context.Context, mod wasm.Module, def plugins.TypeDefinition, data any) (uint32, error) {
	arr, err := utils.ConvertToArray(data)
	if err != nil {
		return 0, err
	}

	var bufferOffset uint32
	var bufferSize uint32

	// write array buffer
	// note: zero-length array has no array buffer
	if len(arr) > 0 {
		itemType := getArraySubtypeInfo(def.Path)
		itemSize := getItemSize(itemType)
		bufferSize = itemSize * uint32(len(arr))
		bufferOffset, err = allocateWasmMemory(ctx, mod, int(bufferSize), 1)
		if err != nil {
			return 0, fmt.Errorf("failed to allocate memory for array buffer: %w", err)
		}

		for i, v := range arr {
			itemOffset := bufferOffset + (itemSize * uint32(i))
			err = writeField(ctx, mod, itemType, itemOffset, v)
			if err != nil {
				return 0, fmt.Errorf("failed to write array item: %w", err)
			}
		}
	}

	// write array object
	offset, err := allocateWasmMemory(ctx, mod, int(def.Size), def.Id)
	if err != nil {
		return 0, err
	}

	mem := mod.Memory()
	ok := mem.WriteUint32Le(offset, bufferOffset)
	if !ok {
		return 0, fmt.Errorf("failed to write array buffer pointer")
	}

	ok = mem.WriteUint32Le(offset+4, bufferOffset)
	if !ok {
		return 0, fmt.Errorf("failed to write array data start pointer")
	}

	ok = mem.WriteUint32Le(offset+8, bufferSize)
	if !ok {
		return 0, fmt.Errorf("failed to write array bytes length")
	}

	ok = mem.WriteUint32Le(offset+12, uint32(len(arr)))
	if !ok {
		return 0, fmt.Errorf("failed to write array length")
	}

	return offset, nil
}
