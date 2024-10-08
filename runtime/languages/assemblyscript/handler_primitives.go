/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/langsupport/primitives"
	"github.com/hypermodeinc/modus/runtime/utils"

	"golang.org/x/exp/constraints"
)

type primitive interface {
	constraints.Integer | constraints.Float | ~bool
}

func (p *planner) NewPrimitiveHandler(ti langsupport.TypeInfo) (h langsupport.TypeHandler, err error) {
	defer func() {
		if err == nil {
			p.typeHandlers[ti.Name()] = h
		}
	}()

	switch ti.Name() {
	case "bool":
		return newPrimitiveHandler[bool](ti), nil
	case "u8":
		return newPrimitiveHandler[uint8](ti), nil
	case "u16":
		return newPrimitiveHandler[uint16](ti), nil
	case "u32":
		return newPrimitiveHandler[uint32](ti), nil
	case "u64":
		return newPrimitiveHandler[uint64](ti), nil
	case "i8":
		return newPrimitiveHandler[int8](ti), nil
	case "i16":
		return newPrimitiveHandler[int16](ti), nil
	case "i32":
		return newPrimitiveHandler[int32](ti), nil
	case "i64":
		return newPrimitiveHandler[int64](ti), nil
	case "f32":
		return newPrimitiveHandler[float32](ti), nil
	case "f64":
		return newPrimitiveHandler[float64](ti), nil
	case "isize":
		return newPrimitiveHandler[int](ti), nil
	case "usize":
		return newPrimitiveHandler[uint](ti), nil
	default:
		return nil, fmt.Errorf("unsupported primitive type: %s", ti.Name())
	}
}

func newPrimitiveHandler[T primitive](ti langsupport.TypeInfo) *primitiveHandler[T] {
	return &primitiveHandler[T]{
		*NewTypeHandler(ti),
		primitives.NewPrimitiveTypeConverter[T](),
	}
}

type primitiveHandler[T primitive] struct {
	typeHandler
	converter primitives.TypeConverter[T]
}

func (h *primitiveHandler[T]) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	val, ok := h.converter.Read(wa.Memory(), offset)
	if !ok {
		return 0, fmt.Errorf("failed to read %s from memory", h.typeInfo.Name())
	}

	return val, nil
}

func (h *primitiveHandler[T]) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	val, err := utils.Cast[T](obj)
	if err != nil {
		return nil, err
	}

	if ok := h.converter.Write(wa.Memory(), offset, val); !ok {
		return nil, fmt.Errorf("failed to write %s to memory", h.typeInfo.Name())
	}

	return nil, nil
}

func (h *primitiveHandler[T]) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if len(vals) != 1 {
		return nil, fmt.Errorf("expected 1 value, got %d", len(vals))
	}

	return h.converter.Decode(vals[0]), nil
}

func (h *primitiveHandler[T]) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	val, err := utils.Cast[T](obj)
	if err != nil {
		return nil, nil, err
	}

	return []uint64{h.converter.Encode(val)}, nil, nil
}
