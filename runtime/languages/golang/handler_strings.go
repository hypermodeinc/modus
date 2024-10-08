/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"unsafe"

	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/spf13/cast"
)

func (p *planner) NewStringHandler(ti langsupport.TypeInfo) (langsupport.TypeHandler, error) {
	handler := &stringHandler{*NewTypeHandler(ti)}
	p.AddHandler(handler)
	return handler, nil
}

type stringHandler struct {
	typeHandler
}

func (h *stringHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return "", nil
	}

	data, size, err := h.readStringHeader(wa, offset)
	if err != nil {
		return "", err
	}

	return h.doReadString(wa, data, size)
}

func (h *stringHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	str, err := cast.ToStringE(obj)
	if err != nil {
		return nil, err
	}

	ptr, cln, err := h.doWriteString(ctx, wa, str)
	if err != nil {
		return cln, err
	}

	data, size, err := h.readStringHeader(wa, ptr)
	if err != nil {
		return cln, err
	}

	return cln, h.writeStringHeader(wa, data, size, offset)
}

func (h *stringHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if len(vals) != 2 {
		return nil, errors.New("expected 2 values when decoding a string")
	}

	data, size := uint32(vals[0]), uint32(vals[1])
	return h.doReadString(wa, data, size)
}

func (h *stringHandler) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	str, err := cast.ToStringE(obj)
	if err != nil {
		return nil, nil, err
	}

	ptr, cln, err := h.doWriteString(ctx, wa, str)
	if err != nil {
		return nil, cln, err
	}

	data, size, err := h.readStringHeader(wa, ptr)
	if err != nil {
		return nil, cln, err
	}

	return []uint64{uint64(data), uint64(size)}, cln, nil
}

func (h *stringHandler) doReadString(wa langsupport.WasmAdapter, offset, size uint32) (string, error) {
	if offset == 0 || size == 0 {
		return "", nil
	}

	bytes, ok := wa.Memory().Read(offset, size)
	if !ok {
		return "", fmt.Errorf("failed to read string data from WASM memory (size: %d)", size)
	}

	return unsafe.String(&bytes[0], size), nil
}

func (h *stringHandler) doWriteString(ctx context.Context, wa langsupport.WasmAdapter, str string) (uint32, utils.Cleaner, error) {
	const id = 2 // ID for string is always 2
	ptr, cln, err := wa.(*wasmAdapter).makeWasmObject(ctx, id, uint32(len(str)))
	if err != nil {
		return 0, cln, err
	}

	offset, ok := wa.Memory().ReadUint32Le(ptr)
	if !ok {
		return 0, cln, errors.New("failed to read string data pointer from WASM memory")
	}

	if ok := wa.Memory().WriteString(offset, str); !ok {
		return 0, cln, fmt.Errorf("failed to write string data to WASM memory (size: %d)", len(str))
	}

	return ptr, cln, nil
}

func (h *stringHandler) readStringHeader(wa langsupport.WasmAdapter, offset uint32) (data, size uint32, err error) {
	if offset == 0 {
		return 0, 0, nil
	}

	val, ok := wa.Memory().ReadUint64Le(offset)
	if !ok {
		return 0, 0, errors.New("failed to read string header from WASM memory")
	}

	data = uint32(val)
	size = uint32(val >> 32)

	return data, size, nil
}

func (h *stringHandler) writeStringHeader(wa langsupport.WasmAdapter, data, size, offset uint32) error {
	val := uint64(size)<<32 | uint64(data)
	if ok := wa.Memory().WriteUint64Le(offset, val); !ok {
		return errors.New("failed to write string header to WASM memory")
	}

	return nil
}
