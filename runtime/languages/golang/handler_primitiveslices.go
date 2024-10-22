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
	"time"

	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/langsupport/primitives"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func (p *planner) NewPrimitiveSliceHandler(ti langsupport.TypeInfo) (h langsupport.TypeHandler, err error) {
	defer func() {
		if err == nil {
			p.typeHandlers[ti.Name()] = h
		}
	}()

	typeDef, err := p.metadata.GetTypeDefinition(ti.Name())
	if err != nil {
		return nil, err
	}

	switch ti.ListElementType().Name() {
	case "bool":
		return newPrimitiveSliceHandler[bool](ti, typeDef), nil
	case "uint8", "byte":
		return newPrimitiveSliceHandler[uint8](ti, typeDef), nil
	case "uint16":
		return newPrimitiveSliceHandler[uint16](ti, typeDef), nil
	case "uint32":
		return newPrimitiveSliceHandler[uint32](ti, typeDef), nil
	case "uint64":
		return newPrimitiveSliceHandler[uint64](ti, typeDef), nil
	case "int8":
		return newPrimitiveSliceHandler[int8](ti, typeDef), nil
	case "int16":
		return newPrimitiveSliceHandler[int16](ti, typeDef), nil
	case "int32", "rune":
		return newPrimitiveSliceHandler[int32](ti, typeDef), nil
	case "int64":
		return newPrimitiveSliceHandler[int64](ti, typeDef), nil
	case "float32":
		return newPrimitiveSliceHandler[float32](ti, typeDef), nil
	case "float64":
		return newPrimitiveSliceHandler[float64](ti, typeDef), nil
	case "int":
		return newPrimitiveSliceHandler[int](ti, typeDef), nil
	case "uint":
		return newPrimitiveSliceHandler[uint](ti, typeDef), nil
	case "uintptr":
		return newPrimitiveSliceHandler[uintptr](ti, typeDef), nil
	case "time.Duration":
		return newPrimitiveSliceHandler[time.Duration](ti, typeDef), nil
	default:
		return nil, fmt.Errorf("unsupported primitive slice type: %s", ti.Name())
	}
}

func newPrimitiveSliceHandler[T primitive](ti langsupport.TypeInfo, typeDef *metadata.TypeDefinition) *primitiveSliceHandler[T] {
	return &primitiveSliceHandler[T]{
		*NewTypeHandler(ti),
		typeDef,
		primitives.NewPrimitiveTypeConverter[T](),
	}
}

type primitiveSliceHandler[T primitive] struct {
	typeHandler
	typeDef   *metadata.TypeDefinition
	converter primitives.TypeConverter[T]
}

func (h *primitiveSliceHandler[T]) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	data, size, _, err := wa.(*wasmAdapter).readSliceHeader(offset)
	if err != nil {
		return nil, err
	}

	return h.doReadSlice(wa, data, size)
}

func (h *primitiveSliceHandler[T]) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := h.doWriteSlice(ctx, wa, obj)
	if err != nil {
		return cln, err
	}

	if ok := utils.CopyMemory(wa.Memory(), ptr, offset, 12); !ok {
		return cln, errors.New("failed to copy slice header")
	}

	return cln, nil
}

func (h *primitiveSliceHandler[T]) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if len(vals) != 3 {
		return nil, errors.New("expected 3 values when decoding a slice")
	}

	// note: capacity is not used here
	data, size := uint32(vals[0]), uint32(vals[1])
	return h.doReadSlice(wa, data, size)
}

func (h *primitiveSliceHandler[T]) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	ptr, cln, err := h.doWriteSlice(ctx, wa, obj)
	if err != nil {
		return nil, cln, err
	}

	data, size, capacity, err := wa.(*wasmAdapter).readSliceHeader(ptr)
	if err != nil {
		return nil, cln, err
	}

	return []uint64{uint64(data), uint64(size), uint64(capacity)}, cln, nil
}

func (h *primitiveSliceHandler[T]) doReadSlice(wa langsupport.WasmAdapter, data, size uint32) (any, error) {
	if data == 0 {
		// nil slice
		return nil, nil
	}

	if size == 0 {
		// empty slice
		return []T{}, nil
	}

	bufferSize := size * uint32(h.converter.TypeSize())
	buf, ok := wa.Memory().Read(data, bufferSize)
	if !ok {
		return nil, errors.New("failed to read data from WASM memory")
	}

	items := h.converter.BytesToSlice(buf)
	return items, nil
}

func (h *primitiveSliceHandler[T]) doWriteSlice(ctx context.Context, wa langsupport.WasmAdapter, obj any) (uint32, utils.Cleaner, error) {
	if utils.HasNil(obj) {
		return 0, nil, nil
	}

	slice, ok := utils.ConvertToSliceOf[T](obj)
	if !ok {
		return 0, nil, fmt.Errorf("expected a %T, got %T", []T{}, obj)
	}

	arrayLen := uint32(len(slice))
	ptr, cln, err := wa.(*wasmAdapter).makeWasmObject(ctx, h.typeDef.Id, arrayLen)
	if err != nil {
		return 0, cln, err
	}

	offset, ok := wa.Memory().ReadUint32Le(ptr)
	if !ok {
		return 0, cln, errors.New("failed to read data pointer from WASM memory")
	}

	bytes := h.converter.SliceToBytes(slice)
	if ok := wa.Memory().Write(offset, bytes); !ok {
		return 0, cln, errors.New("failed to write bytes to WASM memory")
	}

	return ptr, cln, nil
}
