/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"hypruntime/langsupport"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

func (p *planner) NewPointerHandler(ctx context.Context, typ string, rt reflect.Type) (langsupport.TypeHandler, error) {
	handler := NewTypeHandler[pointerHandler](p, typ)
	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)

	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	elementType := _typeInfo.GetUnderlyingType(typ)
	elementHandler, err := p.GetHandler(ctx, elementType)
	if err != nil {
		return nil, err
	}
	handler.elementHandler = elementHandler

	return handler, nil
}

type pointerHandler struct {
	info           langsupport.TypeHandlerInfo
	typeDef        *metadata.TypeDefinition
	elementHandler langsupport.TypeHandler
}

func (h *pointerHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *pointerHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	ptr, ok := wa.Memory().ReadUint32Le(offset)
	if !ok {
		return nil, errors.New("failed to read pointer from memory")
	}

	return h.readData(ctx, wa, ptr)
}

func (h *pointerHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := h.writeData(ctx, wa, obj)
	if err != nil {
		return cln, err
	}

	if ok := wa.Memory().WriteUint32Le(offset, ptr); !ok {
		return cln,
			errors.New("failed to write object pointer to memory")
	}
	return cln, nil
}

func (h *pointerHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if len(vals) != 1 {
		return nil, fmt.Errorf("expected 1 value, got %d", len(vals))
	}

	return h.readData(ctx, wa, uint32(vals[0]))
}

func (h *pointerHandler) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	ptr, cln, err := h.writeData(ctx, wa, obj)
	if err != nil {
		return nil, cln, err
	}

	return []uint64{uint64(ptr)}, cln, nil
}

func (h *pointerHandler) readData(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		// nil pointer
		return nil, nil
	}

	data, err := h.elementHandler.Read(ctx, wa, offset)
	if err != nil {
		return nil, err
	}

	ptr := h.makePointer(data)
	return ptr, nil
}

func (h *pointerHandler) writeData(ctx context.Context, wa langsupport.WasmAdapter, obj any) (uint32, utils.Cleaner, error) {
	if utils.HasNil(obj) {
		// nil pointer
		return 0, nil, nil
	}

	data, err := h.dereferencePointer(obj)
	if err != nil {
		return 0, nil, err
	}

	ptr, cln, err := wa.(*wasmAdapter).newWasmObject(ctx, h.typeDef.Id)
	if err != nil {
		return 0, cln, nil
	}

	c, err := h.elementHandler.Write(ctx, wa, ptr, data)
	cln.AddCleaner(c)
	if err != nil {
		return 0, cln, err
	}

	return ptr, cln, nil
}

func (h *pointerHandler) dereferencePointer(obj any) (any, error) {
	// optimization for common types
	switch t := obj.(type) {
	case *string:
		return *t, nil
	case *bool:
		return *t, nil
	case *int:
		return *t, nil
	case *int8:
		return *t, nil
	case *int16:
		return *t, nil
	case *int32:
		return *t, nil
	case *int64:
		return *t, nil
	case *uint:
		return *t, nil
	case *uint8:
		return *t, nil
	case *uint16:
		return *t, nil
	case *uint32:
		return *t, nil
	case *uint64:
		return *t, nil
	case *uintptr:
		return *t, nil
	case *float32:
		return *t, nil
	case *float64:
		return *t, nil
	case *time.Time:
		return *t, nil
	case *time.Duration:
		return *t, nil
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64,
		time.Time, time.Duration:
		return t, nil // already dereferenced
	}

	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr {
		return obj, nil // already dereferenced
	}
	if rv.IsNil() {
		return nil, errors.New("nil pointer can't be dereferenced")
	}
	return rv.Elem().Interface(), nil
}

func (h *pointerHandler) makePointer(obj any) any {
	// optimization for common types
	switch t := obj.(type) {
	case string:
		return &t
	case bool:
		return &t
	case int:
		return &t
	case int8:
		return &t
	case int16:
		return &t
	case int32:
		return &t
	case int64:
		return &t
	case uint:
		return &t
	case uint8:
		return &t
	case uint16:
		return &t
	case uint32:
		return &t
	case uint64:
		return &t
	case uintptr:
		return &t
	case float32:
		return &t
	case float64:
		return &t
	case time.Time:
		return &t
	case time.Duration:
		return &t
	}

	rt := h.elementHandler.Info().RuntimeType()
	p := reflect.New(rt)
	p.Elem().Set(reflect.ValueOf(obj))
	return p.Interface()
}
