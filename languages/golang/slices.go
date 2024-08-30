/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"hypruntime/utils"
)

func (wa *wasmAdapter) encodeSlice(ctx context.Context, typ string, obj any) ([]uint64, utils.Cleaner, error) {
	ptr, cln, err := wa.doWriteSlice(ctx, typ, obj)
	if err != nil {
		return nil, cln, err
	}

	data, size, capacity, err := wa.readSliceHeader(ptr)
	if err != nil {
		return nil, cln, err
	}

	return []uint64{uint64(data), uint64(size), uint64(capacity)}, cln, nil
}

func (wa *wasmAdapter) decodeSlice(ctx context.Context, typ string, vals []uint64) (any, error) {
	if len(vals) != 3 {
		return nil, errors.New("decodeSlice: expected 3 values")
	}

	// note: capacity is not used here
	data, size := uint32(vals[0]), uint32(vals[1])
	return wa.doReadSlice(ctx, typ, data, size)
}

func (wa *wasmAdapter) writeSlice(ctx context.Context, typ string, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := wa.doWriteSlice(ctx, typ, obj)
	if err != nil {
		return cln, err
	}

	// copy the slice header to the target memory location
	data, size, capacity, err := wa.readSliceHeader(ptr)
	if err != nil {
		return cln, err
	}
	return cln, wa.writeSliceHeader(data, size, capacity, offset)
}

func (wa *wasmAdapter) readSlice(ctx context.Context, typ string, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	data, size, _, err := wa.readSliceHeader(offset)
	if err != nil {
		return nil, err
	}

	return wa.doReadSlice(ctx, typ, data, size)
}

func (wa *wasmAdapter) doReadSlice(ctx context.Context, typ string, data, size uint32) (any, error) {
	if data == 0 {
		// nil slice
		return nil, nil
	}

	if size == 0 {
		// empty slice
		return []any{}, nil
	}

	// special case for byte slices, because they can be read more efficiently
	if typ == "[]byte" || typ == "[]uint8" {
		bytes, ok := wa.mod.Memory().Read(data, size)
		if !ok {
			return nil, errors.New("failed to read data from WASM memory")
		}
		return bytes, nil
	}

	itemType := wa.typeInfo.GetListSubtype(typ)
	itemSize, err := wa.typeInfo.GetSizeOfType(ctx, itemType)
	if err != nil {
		return nil, err
	}

	goItemType, err := wa.getReflectedType(ctx, itemType)
	if err != nil {
		return nil, err
	}

	items := reflect.MakeSlice(reflect.SliceOf(goItemType), int(size), int(size))
	for i := uint32(0); i < size; i++ {
		itemOffset := data + i*itemSize
		item, err := wa.readObject(ctx, itemType, itemOffset)
		if err != nil {
			return nil, err
		}
		items.Index(int(i)).Set(reflect.ValueOf(item))
	}

	return items.Interface(), nil
}

func (wa *wasmAdapter) doWriteSlice(ctx context.Context, typ string, obj any) (uint32, utils.Cleaner, error) {
	if obj == nil {
		return 0, nil, nil
	}

	// special case for byte slices, because they can be written more efficiently
	if typ == "[]byte" || typ == "[]uint8" {
		bytes, err := convertToByteSlice(obj)
		if err != nil {
			return 0, nil, err
		}
		return wa.writeByteSlice(ctx, bytes)
	}

	slice, err := utils.ConvertToSlice(obj)
	if err != nil {
		return 0, nil, err
	}

	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return 0, nil, err
	}

	ptr, cln, err := wa.makeWasmObject(ctx, def.Id, uint32(len(slice)))
	if err != nil {
		return 0, nil, err
	}

	offset, ok := wa.mod.Memory().ReadUint32Le(ptr)
	if !ok {
		return 0, cln, fmt.Errorf("failed to read data pointer from WASM memory")
	}

	t := wa.typeInfo.GetListSubtype(typ)
	size, err := wa.typeInfo.GetSizeOfType(ctx, t)
	if err != nil {
		return 0, cln, err
	}

	for _, val := range slice {
		c, err := wa.writeObject(ctx, t, offset, val)
		cln.AddCleaner(c)
		if err != nil {
			return 0, cln, err
		}
		offset += size
	}

	return ptr, cln, nil
}

func (wa *wasmAdapter) writeSliceHeader(data, size, capacity, offset uint32) error {
	mem := wa.mod.Memory()
	val := uint64(size)<<32 | uint64(data)
	if ok := mem.WriteUint64Le(offset, val); !ok {
		return fmt.Errorf("failed to write slice header to WASM memory")
	}

	if ok := mem.WriteUint32Le(offset+8, capacity); !ok {
		return fmt.Errorf("failed to write slice capacity to WASM memory")
	}

	return nil
}

func (wa *wasmAdapter) readSliceHeader(offset uint32) (data, size, capacity uint32, err error) {
	if offset == 0 {
		return 0, 0, 0, nil
	}

	mem := wa.mod.Memory()
	val, ok := mem.ReadUint64Le(offset)
	if !ok {
		return 0, 0, 0, errors.New("failed to read slice header from WASM memory")
	}

	data = uint32(val)
	size = uint32(val >> 32)

	capacity, ok = mem.ReadUint32Le(offset + 8)
	if !ok {
		return 0, 0, 0, errors.New("failed to read slice capacity from WASM memory")
	}

	return data, size, capacity, nil
}

func (wa *wasmAdapter) writeByteSlice(ctx context.Context, bytes []byte) (uint32, utils.Cleaner, error) {
	const id = 1 // ID for byte slice is always 1
	ptr, cln, err := wa.makeWasmObject(ctx, id, uint32(len(bytes)))
	if err != nil {
		return 0, cln, err
	}

	mem := wa.mod.Memory()
	offset, ok := mem.ReadUint32Le(ptr)
	if !ok {
		return 0, cln, fmt.Errorf("failed to read data pointer from WASM memory")
	}

	if ok := wa.mod.Memory().Write(offset, bytes); !ok {
		return 0, cln, fmt.Errorf("failed to write bytes to WASM memory")
	}

	return ptr, cln, nil
}

func convertToByteSlice(obj any) ([]byte, error) {
	switch obj := obj.(type) {
	case []byte:
		return obj, nil
	case string:
		return []byte(obj), nil
	}

	return nil, fmt.Errorf("input value cannot be used as a byte slice")
}
