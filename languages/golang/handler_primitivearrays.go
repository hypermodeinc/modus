/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"hypruntime/langsupport"
	"hypruntime/langsupport/primitives"
	"hypruntime/utils"
)

func (p *planner) NewPrimitiveArrayHandler(typ string, rt reflect.Type) (langsupport.TypeHandler, error) {
	arrayLen, err := _typeInfo.ArrayLength(typ)
	if err != nil {
		return nil, err
	}

	info := langsupport.NewTypeHandlerInfo(typ, rt, uint32(arrayLen), arrayLen)

	var handler langsupport.TypeHandler

	switch rt.Elem().Kind() {
	case reflect.Bool:
		converter := primitives.NewPrimitiveTypeConverter[bool]()
		handler = &primitiveArrayHandler[bool]{info, converter, arrayLen}
	case reflect.Int8:
		converter := primitives.NewPrimitiveTypeConverter[int8]()
		handler = &primitiveArrayHandler[int8]{info, converter, arrayLen}
	case reflect.Int16:
		converter := primitives.NewPrimitiveTypeConverter[int16]()
		handler = &primitiveArrayHandler[int16]{info, converter, arrayLen}
	case reflect.Int32:
		converter := primitives.NewPrimitiveTypeConverter[int32]()
		handler = &primitiveArrayHandler[int32]{info, converter, arrayLen}
	case reflect.Int64:
		converter := primitives.NewPrimitiveTypeConverter[int64]()
		handler = &primitiveArrayHandler[int64]{info, converter, arrayLen}
	case reflect.Uint8:
		converter := primitives.NewPrimitiveTypeConverter[uint8]()
		handler = &primitiveArrayHandler[uint8]{info, converter, arrayLen}
	case reflect.Uint16:
		converter := primitives.NewPrimitiveTypeConverter[uint16]()
		handler = &primitiveArrayHandler[uint16]{info, converter, arrayLen}
	case reflect.Uint32:
		converter := primitives.NewPrimitiveTypeConverter[uint32]()
		handler = &primitiveArrayHandler[uint32]{info, converter, arrayLen}
	case reflect.Uint64:
		converter := primitives.NewPrimitiveTypeConverter[uint64]()
		handler = &primitiveArrayHandler[uint64]{info, converter, arrayLen}
	case reflect.Float32:
		converter := primitives.NewPrimitiveTypeConverter[float32]()
		handler = &primitiveArrayHandler[float32]{info, converter, arrayLen}
	case reflect.Float64:
		converter := primitives.NewPrimitiveTypeConverter[float64]()
		handler = &primitiveArrayHandler[float64]{info, converter, arrayLen}
	case reflect.Int:
		converter := primitives.NewPrimitiveTypeConverter[int]()
		handler = &primitiveArrayHandler[int]{info, converter, arrayLen}
	case reflect.Uint:
		converter := primitives.NewPrimitiveTypeConverter[uint]()
		handler = &primitiveArrayHandler[uint]{info, converter, arrayLen}
	case reflect.Uintptr:
		converter := primitives.NewPrimitiveTypeConverter[uintptr]()
		handler = &primitiveArrayHandler[uintptr]{info, converter, arrayLen}
	default:
		return nil, fmt.Errorf("unsupported primitive array type: %s", typ)
	}

	p.typeHandlers[typ] = handler
	return handler, nil
}

type primitiveArrayHandler[T primitive] struct {
	info      langsupport.TypeHandlerInfo
	converter primitives.TypeConverter[T]
	arrayLen  int
}

func (h *primitiveArrayHandler[T]) Info() langsupport.TypeHandlerInfo {
	return h.info
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
	array := reflect.New(h.info.RuntimeType()).Elem()
	for i := 0; i < h.arrayLen; i++ {
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
		return nil, fmt.Errorf("failed to write primitive array to WASM memory of type %s", h.info.TypeName())
	}

	return nil, nil
}

func (h *primitiveArrayHandler[T]) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if h.arrayLen == 0 {
		return [0]T{}, nil
	}

	rvItems := reflect.New(h.info.RuntimeType()).Elem()
	for i := 0; i < h.arrayLen; i++ {
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
	for i := 0; i < h.arrayLen; i++ {
		results[i] = h.converter.Encode(items[i])
	}
	return results, nil, nil
}
