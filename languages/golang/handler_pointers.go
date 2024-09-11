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

	ptr := utils.MakePointer(data)
	return ptr, nil
}

func (h *pointerHandler) writeData(ctx context.Context, wa langsupport.WasmAdapter, obj any) (uint32, utils.Cleaner, error) {
	if utils.HasNil(obj) {
		// nil pointer
		return 0, nil, nil
	}

	data := utils.DereferencePointer(obj)

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
