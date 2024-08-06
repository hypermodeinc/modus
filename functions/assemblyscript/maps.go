/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"reflect"

	"hmruntime/functions/assemblyscript/hash"
	"hmruntime/plugins/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

// Reference: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/map.ts

func readMap(ctx context.Context, mem wasm.Memory, def *metadata.TypeDefinition, offset uint32) (data any, err error) {

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

	// the length of array buffer is stored 4 bytes before the offset
	byteLength, ok := mem.ReadUint32Le(entries - 4)
	if !ok {
		return nil, fmt.Errorf("failed to read map entries buffer length")
	}

	entrySize := byteLength / entriesCapacity
	keyType, valueType := getMapSubtypeInfo(def.Path)
	valueOffset := getSizeForOffset(keyType.Path)

	goKeyType, err := getGoType(keyType.Path)
	if err != nil {
		return nil, err
	}
	goValueType, err := getGoType(valueType.Path)
	if err != nil {
		return nil, err
	}

	m := reflect.MakeMapWithSize(reflect.MapOf(goKeyType, goValueType), int(entriesCount))
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
		m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	}

	return m.Interface(), nil
}

func writeMap(ctx context.Context, mod wasm.Module, def *metadata.TypeDefinition, data any) (offset uint32, err error) {

	// Unfortunately, there's no way to do this without reflection.
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Map {
		// TODO: support []kvp ?
		return 0, fmt.Errorf("unsupported map type %T", data)
	}
	mapLen := uint32(rv.Len())

	// unpin everything when done
	var pins = make([]uint32, 0, (mapLen*2)+2)
	defer func() {
		for _, ptr := range pins {
			err = unpinWasmMemory(ctx, mod, ptr)
			if err != nil {
				break
			}
		}
	}()

	// determine capacities and mask
	bucketsCapacity := uint32(4)
	entriesCapacity := uint32(4)
	bucketsMask := bucketsCapacity - 1
	for bucketsCapacity < mapLen {
		bucketsCapacity <<= 1
		entriesCapacity = bucketsCapacity * 8 / 3
		bucketsMask = bucketsCapacity - 1
	}

	// create buckets array buffer
	const bucketSize = 4
	bucketsBufferSize := bucketSize * bucketsCapacity
	bucketsBufferOffset, err := allocateWasmMemory(ctx, mod, bucketsBufferSize, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to allocate memory for array buffer: %w", err)
	}

	// pin the array buffer so it can't get garbage collected
	// when we allocate the array object
	err = pinWasmMemory(ctx, mod, bucketsBufferOffset)
	if err != nil {
		return 0, fmt.Errorf("failed to pin array buffer: %w", err)
	}
	pins = append(pins, bucketsBufferOffset)

	// write entries array buffer
	// note: unlike arrays, an empty map DOES have array buffers
	keyType, valueType := getMapSubtypeInfo(def.Path)
	keySize := getItemSize(keyType)
	valueSize := getItemSize(valueType)
	const taggedNextSize = 4
	entryAlign := max(keySize, valueSize, 4) - 1
	entrySize := (keySize + valueSize + taggedNextSize + entryAlign) & ^entryAlign
	entriesBufferSize := entrySize * entriesCapacity
	entriesBufferOffset, err := allocateWasmMemory(ctx, mod, entriesBufferSize, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to allocate memory for array buffer: %w", err)
	}

	// pin the array buffer so it can't get garbage collected
	// when we allocate the array object
	err = pinWasmMemory(ctx, mod, entriesBufferOffset)
	if err != nil {
		return 0, fmt.Errorf("failed to pin array buffer: %w", err)
	}
	pins = append(pins, entriesBufferOffset)

	valueOffset := getSizeForOffset(keyType.Path)
	taggedNextOffset := getSizeForOffset(valueType.Path) + valueOffset

	mem := mod.Memory()
	mapKeys := rv.MapKeys()
	for i, mapKey := range mapKeys {

		entryOffset := entriesBufferOffset + (entrySize * uint32(i))

		// write entry key and calculate hash code
		var hashCode, ptr uint32
		key := mapKey.Interface()
		switch t := key.(type) {
		case string:
			// Special case for string keys.  Since we need to encode them as UTF16,
			// for both writing to memory and calculating the hash code, bypass the
			// normal writeField/writeString functions and do it manually.
			bytes := encodeUTF16(t)
			hashCode = hash.GetHashCode(bytes)
			ptr, err = writeRawBytes(ctx, mod, bytes, 2)
			if err != nil {
				return 0, fmt.Errorf("failed to write map entry key: %w", err)
			}
			ok := mem.WriteUint32Le(entryOffset, ptr)
			if !ok {
				return 0, fmt.Errorf("failed to write map entry key pointer")
			}

		default:
			hashCode = hash.GetHashCode(key)
			ptr, err = writeField(ctx, mod, keyType, entryOffset, key)
			if err != nil {
				return 0, fmt.Errorf("failed to write map entry key: %w", err)
			}
		}

		// If we allocated memory for the key, we need to pin it too.
		if ptr != 0 {
			err = pinWasmMemory(ctx, mod, ptr)
			if err != nil {
				return 0, err
			}
			pins = append(pins, ptr)
		}

		// write entry value
		mapValue := rv.MapIndex(mapKey)
		value := mapValue.Interface()
		entryValueOffset := entryOffset + valueOffset
		ptr, err = writeField(ctx, mod, valueType, entryValueOffset, value)
		if err != nil {
			return 0, fmt.Errorf("failed to write map entry value: %w", err)
		}

		// If we allocated memory for the value, we need to pin it too.
		if ptr != 0 {
			err = pinWasmMemory(ctx, mod, ptr)
			if err != nil {
				return 0, err
			}
			pins = append(pins, ptr)
		}

		// write to bucket and "tagged next" field
		bucketPtrBase := bucketsBufferOffset + ((hashCode & bucketsMask) * bucketSize)
		prev, ok := mem.ReadUint32Le(bucketPtrBase)
		if !ok {
			return 0, fmt.Errorf("failed to read previous map entry bucket pointer")
		}
		ok = mem.WriteUint32Le(entryOffset+taggedNextOffset, prev)
		if !ok {
			return 0, fmt.Errorf("failed to write map entry tagged next field")
		}
		ok = mem.WriteUint32Le(bucketPtrBase, entryOffset)
		if !ok {
			return 0, fmt.Errorf("failed to write map entry bucket pointer")
		}
	}

	// write map object
	const size = 24
	offset, err = allocateWasmMemory(ctx, mod, size, def.Id)
	if err != nil {
		return 0, err
	}

	ok := mem.WriteUint32Le(offset, bucketsBufferOffset)
	if !ok {
		return 0, fmt.Errorf("failed to write map buckets pointer")
	}

	ok = mem.WriteUint32Le(offset+4, bucketsMask)
	if !ok {
		return 0, fmt.Errorf("failed to write map buckets mask")
	}

	ok = mem.WriteUint32Le(offset+8, entriesBufferOffset)
	if !ok {
		return 0, fmt.Errorf("failed to write map entries pointer")
	}

	ok = mem.WriteUint32Le(offset+12, entriesCapacity)
	if !ok {
		return 0, fmt.Errorf("failed to write map entries capacity")
	}

	ok = mem.WriteUint32Le(offset+16, mapLen)
	if !ok {
		return 0, fmt.Errorf("failed to write map entries offset")
	}

	ok = mem.WriteUint32Le(offset+20, mapLen)
	if !ok {
		return 0, fmt.Errorf("failed to write map entries count")
	}

	return offset, nil
}

func getSizeForOffset(t string) uint32 {
	switch t {
	case "u64", "i64", "f64":
		// 64-bit keys have 8-byte value offset
		return 8
	default:
		// everything else has 4-byte value offset
		return 4
	}
}
