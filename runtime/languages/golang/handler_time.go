/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang

import (
	"context"
	"errors"
	"time"
	"unsafe"

	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func (p *planner) NewTimeHandler(ti langsupport.TypeInfo) (langsupport.TypeHandler, error) {
	handler := &timeHandler{*NewTypeHandler(ti)}
	p.AddHandler(handler)
	return handler, nil
}

type timeHandler struct {
	typeHandler
}

func (h *timeHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	wall, ok := wa.Memory().ReadUint64Le(offset)
	if !ok {
		return nil, errors.New("failed to read time.Time.wall from WASM memory")
	}

	x, ok := wa.Memory().ReadUint64Le(offset + 8)
	if !ok {
		return nil, errors.New("failed to read time.Time.ext from WASM memory")
	}
	ext := int64(x)

	// skip loc - we only support UTC

	return timeFromVals(wall, ext), nil
}

func (h *timeHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	tm, err := utils.ConvertToTimestamp(obj)
	if err != nil {
		return nil, err
	}

	wall, ext := getTimeVals(tm)

	if !wa.Memory().WriteUint64Le(offset, wall) {
		return nil, errors.New("failed to write time.Time.wall to WASM memory")
	}

	if !wa.Memory().WriteUint64Le(offset+8, uint64(ext)) {
		return nil, errors.New("failed to write time.Time.ext to WASM memory")
	}

	// skip loc - we only support UTC

	return nil, nil
}

func (h *timeHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if len(vals) != 3 {
		return nil, errors.New("decodeTime: expected 3 values")
	}

	wall, ext := vals[0], int64(vals[1])
	// skip loc - we only support UTC

	return timeFromVals(wall, ext), nil
}

func (h *timeHandler) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	tm, err := utils.ConvertToTimestamp(obj)
	if err != nil {
		return []uint64{0}, nil, err
	}

	wall, ext := getTimeVals(tm)

	// skip loc - we only support UTC

	return []uint64{wall, uint64(ext), 0}, nil, nil
}

func timeFromVals(wall uint64, ext int64) time.Time {
	type tm struct {
		wall uint64
		ext  int64
		loc  *time.Location
	}

	t := tm{wall, ext, nil}
	return *(*time.Time)(unsafe.Pointer(&t))
}

func getTimeVals(t time.Time) (uint64, int64) {
	type tm struct {
		wall uint64
		ext  int64
		loc  *time.Location
	}

	s := *(*tm)(unsafe.Pointer(&t))
	return s.wall, s.ext
}
