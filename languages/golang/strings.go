/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"unsafe"

	"hypruntime/utils"

	"github.com/spf13/cast"
)

func (wa *wasmAdapter) encodeString(ctx context.Context, obj any) ([]uint64, utils.Cleaner, error) {
	str, err := cast.ToStringE(obj)
	if err != nil {
		return nil, nil, err
	}

	ptr, cln, err := wa.doWriteString(ctx, str)
	if err != nil {
		return nil, cln, err
	}

	data, size, err := wa.readStringHeader(ptr)
	if err != nil {
		return nil, cln, err
	}

	return []uint64{uint64(data), uint64(size)}, cln, nil
}

func (wa *wasmAdapter) decodeString(vals []uint64) (any, error) {
	if len(vals) != 2 {
		return nil, errors.New("decodeString: expected 2 values")
	}

	data, size := uint32(vals[0]), uint32(vals[1])
	return wa.doReadString(data, size)
}

func (wa *wasmAdapter) writeString(ctx context.Context, offset uint32, obj any) (utils.Cleaner, error) {
	str, err := cast.ToStringE(obj)
	if err != nil {
		return nil, err
	}
	ptr, cln, err := wa.doWriteString(ctx, str)
	if err != nil {
		return cln, err
	}
	data, size, err := wa.readStringHeader(ptr)
	if err != nil {
		return cln, err
	}
	return cln, wa.writeStringHeader(data, size, offset)
}

func (wa *wasmAdapter) readString(offset uint32) (string, error) {
	if offset == 0 {
		return "", nil
	}

	data, size, err := wa.readStringHeader(offset)
	if err != nil {
		return "", err
	}

	return wa.doReadString(data, size)
}

func (wa *wasmAdapter) doReadString(offset, size uint32) (string, error) {
	if offset == 0 || size == 0 {
		return "", nil
	}

	bytes, ok := wa.mod.Memory().Read(offset, size)
	if !ok {
		return "", errors.New("failed to read data from WASM memory")
	}

	return unsafe.String(&bytes[0], size), nil
}

func (wa *wasmAdapter) doWriteString(ctx context.Context, s string) (uint32, utils.Cleaner, error) {
	const id = 2 // ID for string is always 2
	ptr, cln, err := wa.makeWasmObject(ctx, id, uint32(len(s)))
	if err != nil {
		return 0, cln, err
	}

	mem := wa.mod.Memory()
	offset, ok := mem.ReadUint32Le(ptr)
	if !ok {
		return 0, cln, fmt.Errorf("failed to read data pointer from WASM memory")
	}

	if ok := wa.mod.Memory().WriteString(offset, s); !ok {
		return 0, cln, fmt.Errorf("failed to write string to WASM memory")
	}

	return ptr, cln, nil
}

func (wa *wasmAdapter) writeStringHeader(data, size, offset uint32) error {
	val := uint64(size)<<32 | uint64(data)
	if ok := wa.mod.Memory().WriteUint64Le(offset, val); !ok {
		return fmt.Errorf("failed to write string header to WASM memory")
	}
	return nil
}

func (wa *wasmAdapter) readStringHeader(offset uint32) (data, size uint32, err error) {
	if offset == 0 {
		return 0, 0, nil
	}

	val, ok := wa.mod.Memory().ReadUint64Le(offset)
	if !ok {
		return 0, 0, errors.New("failed to read string header from WASM memory")
	}

	data = uint32(val)
	size = uint32(val >> 32)

	return data, size, nil
}
