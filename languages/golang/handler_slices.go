/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"reflect"

	"hypruntime/langsupport"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

func (p *planner) NewSliceHandler(ctx context.Context, typ string, rt reflect.Type) (langsupport.TypeHandler, error) {
	handler := NewTypeHandler[sliceHandler](p, typ)

	// slice header is a 4 byte pointer, 4 byte length, 4 byte capacity
	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, 12, 3)

	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	elementType := _typeInfo.GetListSubtype(typ)
	elementHandler, err := p.GetHandler(ctx, elementType)
	if err != nil {
		return nil, err
	}
	handler.elementHandler = elementHandler

	handler.emptyValue = reflect.MakeSlice(rt, 0, 0).Interface()

	return handler, nil
}

type sliceHandler struct {
	info           langsupport.TypeHandlerInfo
	typeDef        *metadata.TypeDefinition
	elementHandler langsupport.TypeHandler
	emptyValue     any
}

func (h *sliceHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *sliceHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	data, size, _, err := wa.(*wasmAdapter).readSliceHeader(offset)
	if err != nil {
		return nil, err
	}

	return h.doReadSlice(ctx, wa, data, size)
}

func (h *sliceHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := h.doWriteSlice(ctx, wa, obj)
	if err != nil {
		return cln, err
	}

	if ok := utils.CopyMemory(wa.Memory(), ptr, offset, 12); !ok {
		return cln, errors.New("failed to copy slice header")
	}

	return cln, nil
}

func (h *sliceHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if len(vals) != 3 {
		return nil, errors.New("expected 3 values when decoding a slice")
	}

	// note: capacity is not used here
	data, size := uint32(vals[0]), uint32(vals[1])
	return h.doReadSlice(ctx, wa, data, size)
}

func (h *sliceHandler) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
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

func (h *sliceHandler) doReadSlice(ctx context.Context, wa langsupport.WasmAdapter, data, size uint32) (any, error) {
	if data == 0 {
		// nil slice
		return nil, nil
	}

	if size == 0 {
		// empty slice
		return h.emptyValue, nil
	}

	elementSize := h.elementHandler.Info().TypeSize()
	items := reflect.MakeSlice(h.Info().RuntimeType(), int(size), int(size))
	for i := uint32(0); i < size; i++ {
		itemOffset := data + i*elementSize
		item, err := h.elementHandler.Read(ctx, wa, itemOffset)
		if err != nil {
			return nil, err
		}
		items.Index(int(i)).Set(reflect.ValueOf(item))
	}

	return items.Interface(), nil
}

func (h *sliceHandler) doWriteSlice(ctx context.Context, wa langsupport.WasmAdapter, obj any) (ptr uint32, cln utils.Cleaner, err error) {
	if utils.HasNil(obj) {
		return 0, nil, nil
	}

	slice, err := utils.ConvertToSlice(obj)
	if err != nil {
		return 0, nil, err
	}

	ptr, cln, err = wa.(*wasmAdapter).makeWasmObject(ctx, h.typeDef.Id, uint32(len(slice)))
	if err != nil {
		return 0, nil, err
	}

	offset, ok := wa.Memory().ReadUint32Le(ptr)
	if !ok {
		return 0, cln, errors.New("failed to read data pointer from WASM memory")
	}

	innerCln := utils.NewCleanerN(len(slice))
	defer func() {
		// unpin slice elements after the slice is written to memory
		if e := innerCln.Clean(); e != nil && err == nil {
			err = e
		}
	}()

	elementSize := h.elementHandler.Info().TypeSize()
	for _, val := range slice {
		c, err := h.elementHandler.Write(ctx, wa, offset, val)
		innerCln.AddCleaner(c)
		if err != nil {
			return 0, cln, err
		}
		offset += elementSize
	}

	return ptr, cln, nil
}

func (wa *wasmAdapter) readSliceHeader(offset uint32) (data, size, capacity uint32, err error) {
	if offset == 0 {
		return 0, 0, 0, nil
	}

	val, ok := wa.Memory().ReadUint64Le(offset)
	if !ok {
		return 0, 0, 0, errors.New("failed to read slice header from WASM memory")
	}

	data = uint32(val)
	size = uint32(val >> 32)

	capacity, ok = wa.Memory().ReadUint32Le(offset + 8)
	if !ok {
		return 0, 0, 0, errors.New("failed to read slice capacity from WASM memory")
	}

	return data, size, capacity, nil
}