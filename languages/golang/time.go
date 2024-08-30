/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unsafe"

	"hypruntime/utils"
)

func (wa *wasmAdapter) encodeTime(ctx context.Context, obj any) ([]uint64, utils.Cleaner, error) {
	offset, cln, err := wa.doWriteTime(ctx, obj)
	if err != nil {
		return nil, cln, err
	}

	mem := wa.mod.Memory()
	wall, ok := mem.ReadUint64Le(offset)
	if !ok {
		return nil, cln, errors.New("failed to read time.Time.wall from WASM memory")
	}

	ext, ok := mem.ReadUint64Le(offset + 8)
	if !ok {
		return nil, cln, errors.New("failed to read time.Time.ext from WASM memory")
	}

	// skip loc - we only support UTC

	return []uint64{wall, ext, 0}, cln, nil
}

func (wa *wasmAdapter) decodeTime(vals []uint64) (any, error) {
	if len(vals) != 3 {
		return nil, errors.New("decodeTime: expected 3 values")
	}

	wall, ext := vals[0], int64(vals[1])
	// skip loc - we only support UTC

	return timeFromVals(wall, ext), nil
}

func (wa *wasmAdapter) writeTime(ctx context.Context, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := wa.doWriteTime(ctx, obj)
	if err != nil {
		return cln, err
	}

	if !wa.mod.Memory().WriteUint32Le(offset, ptr) {
		return cln, errors.New("failed to write time.Time internal pointer to memory")
	}

	return cln, nil
}

func (wa *wasmAdapter) readTime(offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	mem := wa.mod.Memory()
	wall, ok := mem.ReadUint64Le(offset)
	if !ok {
		return nil, errors.New("failed to read time.Time.wall from WASM memory")
	}

	x, ok := mem.ReadUint64Le(offset + 8)
	if !ok {
		return nil, errors.New("failed to read time.Time.ext from WASM memory")
	}
	ext := int64(x)

	// skip loc - we only support UTC

	return timeFromVals(wall, ext), nil
}

func (wa *wasmAdapter) doWriteTime(ctx context.Context, obj any) (uint32, utils.Cleaner, error) {
	t, ok := obj.(time.Time)
	if !ok {
		return 0, nil, fmt.Errorf("expected time.Time, got %T", obj)
	}

	def, err := wa.typeInfo.GetTypeDefinition(ctx, "time.Time")
	if err != nil {
		return 0, nil, err
	}

	ptr, cln, err := wa.newWasmObject(ctx, def.Id)
	if err != nil {
		return 0, cln, nil
	}

	wall, ext := getTimeVals(t)

	mem := wa.mod.Memory()
	if !mem.WriteUint64Le(ptr, wall) {
		return 0, cln, errors.New("failed to write time.Time.wall to WASM memory")
	}

	if !mem.WriteUint64Le(ptr+8, uint64(ext)) {
		return 0, cln, errors.New("failed to write time.Time.ext to WASM memory")
	}

	// skip loc - we only support UTC

	return ptr, cln, nil
}

type tm struct {
	wall uint64
	ext  int64
	loc  *time.Location
}

func timeFromVals(wall uint64, ext int64) time.Time {
	t := tm{wall, ext, nil}
	return *(*time.Time)(unsafe.Pointer(&t))
}

func getTimeVals(t time.Time) (uint64, int64) {
	tm := *(*tm)(unsafe.Pointer(&t))
	return tm.wall, tm.ext
}
