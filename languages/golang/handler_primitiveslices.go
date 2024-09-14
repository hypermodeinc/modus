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
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

func (p *planner) NewPrimitiveSliceHandler(typ string, rt reflect.Type) (langsupport.TypeHandler, error) {
	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}

	// slice header is a 4 byte pointer, 4 byte length, 4 byte capacity
	info := langsupport.NewTypeHandlerInfo(typ, rt, 12, 3)

	var handler langsupport.TypeHandler

	switch rt.Elem().Kind() {
	case reflect.Bool:
		converter := primitives.NewPrimitiveTypeConverter[bool]()
		handler = &primitiveSliceHandler[bool]{info, typeDef, converter}
	case reflect.Int8:
		converter := primitives.NewPrimitiveTypeConverter[int8]()
		handler = &primitiveSliceHandler[int8]{info, typeDef, converter}
	case reflect.Int16:
		converter := primitives.NewPrimitiveTypeConverter[int16]()
		handler = &primitiveSliceHandler[int16]{info, typeDef, converter}
	case reflect.Int32:
		converter := primitives.NewPrimitiveTypeConverter[int32]()
		handler = &primitiveSliceHandler[int32]{info, typeDef, converter}
	case reflect.Int64:
		converter := primitives.NewPrimitiveTypeConverter[int64]()
		handler = &primitiveSliceHandler[int64]{info, typeDef, converter}
	case reflect.Uint8:
		converter := primitives.NewPrimitiveTypeConverter[uint8]()
		handler = &primitiveSliceHandler[uint8]{info, typeDef, converter}
	case reflect.Uint16:
		converter := primitives.NewPrimitiveTypeConverter[uint16]()
		handler = &primitiveSliceHandler[uint16]{info, typeDef, converter}
	case reflect.Uint32:
		converter := primitives.NewPrimitiveTypeConverter[uint32]()
		handler = &primitiveSliceHandler[uint32]{info, typeDef, converter}
	case reflect.Uint64:
		converter := primitives.NewPrimitiveTypeConverter[uint64]()
		handler = &primitiveSliceHandler[uint64]{info, typeDef, converter}
	case reflect.Float32:
		converter := primitives.NewPrimitiveTypeConverter[float32]()
		handler = &primitiveSliceHandler[float32]{info, typeDef, converter}
	case reflect.Float64:
		converter := primitives.NewPrimitiveTypeConverter[float64]()
		handler = &primitiveSliceHandler[float64]{info, typeDef, converter}
	case reflect.Int:
		converter := primitives.NewPrimitiveTypeConverter[int]()
		handler = &primitiveSliceHandler[int]{info, typeDef, converter}
	case reflect.Uint:
		converter := primitives.NewPrimitiveTypeConverter[uint]()
		handler = &primitiveSliceHandler[uint]{info, typeDef, converter}
	case reflect.Uintptr:
		converter := primitives.NewPrimitiveTypeConverter[uintptr]()
		handler = &primitiveSliceHandler[uintptr]{info, typeDef, converter}
	default:
		return nil, fmt.Errorf("unsupported primitive slice type: %s", typ)
	}

	p.typeHandlers[typ] = handler
	return handler, nil
}

type primitiveSliceHandler[T primitive] struct {
	info      langsupport.TypeHandlerInfo
	typeDef   *metadata.TypeDefinition
	converter primitives.TypeConverter[T]
}

func (h *primitiveSliceHandler[T]) Info() langsupport.TypeHandlerInfo {
	return h.info
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

	items, ok := obj.([]T)
	if !ok {
		return 0, nil, fmt.Errorf("expected a %T, got %T", []T{}, obj)
	}

	arrayLen := uint32(len(items))
	ptr, cln, err := wa.(*wasmAdapter).makeWasmObject(ctx, h.typeDef.Id, arrayLen)
	if err != nil {
		return 0, cln, err
	}

	offset, ok := wa.Memory().ReadUint32Le(ptr)
	if !ok {
		return 0, cln, errors.New("failed to read data pointer from WASM memory")
	}

	bytes := h.converter.SliceToBytes(items)
	if ok := wa.Memory().Write(offset, bytes); !ok {
		return 0, cln, errors.New("failed to write bytes to WASM memory")
	}

	return ptr, cln, nil
}
