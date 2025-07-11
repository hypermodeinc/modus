/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/langsupport/primitives"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func (p *planner) NewPrimitiveArrayHandler(ti langsupport.TypeInfo) (h langsupport.TypeHandler, err error) {
	defer func() {
		if err == nil {
			p.typeHandlers[ti.Name()] = h
		}
	}()

	switch ti.ListElementType().Name() {
	case "bool":
		return newPrimitiveArrayHandler[bool](ti)
	case "uint8", "byte":
		return newPrimitiveArrayHandler[uint8](ti)
	case "uint16":
		return newPrimitiveArrayHandler[uint16](ti)
	case "uint32":
		return newPrimitiveArrayHandler[uint32](ti)
	case "uint64":
		return newPrimitiveArrayHandler[uint64](ti)
	case "int8":
		return newPrimitiveArrayHandler[int8](ti)
	case "int16":
		return newPrimitiveArrayHandler[int16](ti)
	case "int32", "rune":
		return newPrimitiveArrayHandler[int32](ti)
	case "int64":
		return newPrimitiveArrayHandler[int64](ti)
	case "float32":
		return newPrimitiveArrayHandler[float32](ti)
	case "float64":
		return newPrimitiveArrayHandler[float64](ti)
	case "int":
		return newPrimitiveArrayHandler[int](ti)
	case "uint":
		return newPrimitiveArrayHandler[uint](ti)
	case "uintptr":
		return newPrimitiveArrayHandler[uintptr](ti)
	case "time.Duration":
		return newPrimitiveArrayHandler[time.Duration](ti)
	default:
		return nil, fmt.Errorf("unsupported primitive array type: %s", ti.Name())
	}
}

func newPrimitiveArrayHandler[T primitive](ti langsupport.TypeInfo) (*primitiveArrayHandler[T], error) {
	arrayLen, err := _langTypeInfo.ArrayLength(ti.Name())
	if err != nil {
		return nil, err
	}

	return &primitiveArrayHandler[T]{
		*NewTypeHandler(ti),
		primitives.NewPrimitiveTypeConverter[T](),
		arrayLen,
	}, nil
}

type primitiveArrayHandler[T primitive] struct {
	typeHandler
	converter primitives.TypeConverter[T]
	arrayLen  int
}

func (h *primitiveArrayHandler[T]) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if h.arrayLen == 0 {
		return [0]T{}, nil
	}

	bufferSize := uint32(h.arrayLen * h.converter.TypeSize())
	buf, ok := wa.Memory().Read(offset, bufferSize)
	if !ok {
		return nil, errors.New("failed to read data from WASM memory")
	}

	items := h.converter.BytesToSlice(buf)

	// convert the slice to an array of the target type
	array := reflect.New(h.typeInfo.ReflectedType()).Elem()
	for i := range h.arrayLen {
		array.Index(i).Set(reflect.ValueOf(items[i]))
	}
	return array.Interface(), nil
}

func (h *primitiveArrayHandler[T]) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	items, ok := utils.ConvertToSliceOf[T](obj)
	if !ok {
		return nil, fmt.Errorf("expected a type compatible with %T, got %T", []T{}, obj)
	}

	// the array must be exactly the right length
	if len(items) < h.arrayLen {
		items = append(items, make([]T, h.arrayLen-len(items))...)
	} else if len(items) > h.arrayLen {
		items = items[:h.arrayLen]
	}

	bytes := h.converter.SliceToBytes(items)
	if ok := wa.Memory().Write(offset, bytes); !ok {
		return nil, fmt.Errorf("failed to write primitive array to WASM memory of type %s", h.typeInfo.Name())
	}

	return nil, nil
}

func (h *primitiveArrayHandler[T]) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if h.arrayLen == 0 {
		return [0]T{}, nil
	}

	rvItems := reflect.New(h.typeInfo.ReflectedType()).Elem()
	for i := range h.arrayLen {
		item := h.converter.Decode(vals[i])
		ptr := rvItems.Index(i).Addr().UnsafePointer()
		*(*T)(ptr) = item
	}
	return rvItems.Interface(), nil
}

func (h *primitiveArrayHandler[T]) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	if h.arrayLen == 0 {
		return nil, nil, nil
	}

	items, ok := utils.ConvertToSliceOf[T](obj)
	if !ok {
		return nil, nil, fmt.Errorf("expected a type compatible with %T, got %T", []T{}, obj)
	}

	results := make([]uint64, h.arrayLen)
	for i := range h.arrayLen {
		results[i] = h.converter.Encode(items[i])
	}
	return results, nil, nil
}
