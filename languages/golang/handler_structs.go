/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"hypruntime/langsupport"
	"hypruntime/logger"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

const maxDepth = 5 // TODO: make this based on the depth requested in the query

func (p *planner) NewStructHandler(ctx context.Context, typ string, rt reflect.Type) (langsupport.TypeHandler, error) {
	handler := NewTypeHandler[structHandler](p, typ)

	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	offset := uint32(0)
	encodingLength := 0
	fieldHandlers := make([]langsupport.TypeHandler, len(typeDef.Fields))
	for i, field := range typeDef.Fields {
		fieldHandler, err := p.GetHandler(ctx, field.Type)
		if err != nil {
			return nil, err
		}
		fieldHandlers[i] = fieldHandler
		size := fieldHandler.Info().TypeSize()
		pad := langsupport.GetAlignmentPadding(offset, size)
		offset += size + pad
		encodingLength += fieldHandler.Info().EncodingLength()
	}
	handler.fieldHandlers = fieldHandlers

	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, offset, encodingLength)

	return handler, nil
}

type structHandler struct {
	info          langsupport.TypeHandlerInfo
	typeDef       *metadata.TypeDefinition
	fieldHandlers []langsupport.TypeHandler
}

func (h *structHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *structHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {

	// Check for recursion
	visitedPtrs := wa.(*wasmAdapter).visitedPtrs
	if visitedPtrs[offset] >= maxDepth {
		logger.Warn(ctx).Bool("user_visible", true).Msgf("Excessive recursion detected in %s. Stopping at depth %d.", h.info.TypeName(), maxDepth)
		return nil, nil
	}
	visitedPtrs[offset]++
	defer func() {
		n := visitedPtrs[offset]
		if n == 1 {
			delete(visitedPtrs, offset)
		} else {
			visitedPtrs[offset] = n - 1
		}
	}()

	m := make(map[string]any, len(h.fieldHandlers))
	fieldOffset := uint32(0)
	for i, field := range h.typeDef.Fields {
		handler := h.fieldHandlers[i]
		size := handler.Info().TypeSize()
		fieldOffset += langsupport.GetAlignmentPadding(fieldOffset, size)

		val, err := handler.Read(ctx, wa, offset+fieldOffset)
		if err != nil {
			return nil, err
		}

		m[field.Name] = val
		fieldOffset += size
	}

	return h.getStructOutput(m)
}

func (h *structHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	var mapObj map[string]any
	var rvObj reflect.Value
	if m, ok := obj.(map[string]any); ok {
		mapObj = m
	} else {
		rvObj = reflect.ValueOf(obj)
		if rvObj.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected a struct, got %s", rvObj.Kind())
		}
	}

	numFields := len(h.typeDef.Fields)
	cleaner := utils.NewCleanerN(numFields)

	fieldOffset := uint32(0)
	for i, field := range h.typeDef.Fields {
		var fieldObj any
		if mapObj != nil {
			// case sensitive when reading from map
			fieldObj = mapObj[field.Name]
		} else {
			// case insensitive when reading from struct
			fieldObj = rvObj.FieldByNameFunc(func(s string) bool { return strings.EqualFold(s, field.Name) }).Interface()
		}

		size := h.fieldHandlers[i].Info().TypeSize()
		pad := langsupport.GetAlignmentPadding(fieldOffset, size)
		fieldOffset += pad

		handler := h.fieldHandlers[i]
		cln, err := handler.Write(ctx, wa, offset+fieldOffset, fieldObj)
		cleaner.AddCleaner(cln)
		if err != nil {
			return cleaner, err
		}

		fieldOffset += size
	}

	return cleaner, nil
}

func (h *structHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	switch len(h.fieldHandlers) {
	case 0:
		return nil, nil
	case 1:
		data, err := h.fieldHandlers[0].Decode(ctx, wa, vals)
		if err != nil {
			return nil, err
		}
		m := map[string]any{h.typeDef.Fields[0].Name: data}
		return h.getStructOutput(m)
	}

	// this doesn't need to be supported until TinyGo implements multi-value returns
	return nil, fmt.Errorf("decoding struct of type %s is not supported", h.info.TypeName())
}

func (h *structHandler) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	var mapObj map[string]any
	var rvObj reflect.Value
	if m, ok := obj.(map[string]any); ok {
		mapObj = m
	} else {
		rvObj = reflect.ValueOf(obj)
		if rvObj.Kind() != reflect.Struct {
			return nil, nil, fmt.Errorf("expected a struct, got %s", rvObj.Kind())
		}
	}

	numFields := len(h.typeDef.Fields)
	results := make([]uint64, 0, numFields*2)
	cleaner := utils.NewCleanerN(numFields)

	for i, field := range h.typeDef.Fields {
		var fieldObj any
		if mapObj != nil {
			// case sensitive when reading from map
			fieldObj = mapObj[field.Name]
		} else {
			// case insensitive when reading from struct
			fieldObj = rvObj.FieldByNameFunc(func(s string) bool { return strings.EqualFold(s, field.Name) }).Interface()
		}

		handler := h.fieldHandlers[i]
		vals, cln, err := handler.Encode(ctx, wa, fieldObj)
		cleaner.AddCleaner(cln)
		if err != nil {
			return nil, cleaner, err
		}
		results = append(results, vals...)
	}

	return results, cleaner, nil
}

func (h structHandler) getStructOutput(data map[string]any) (any, error) {
	rt := h.info.RuntimeType()
	if rt.Kind() == reflect.Map {
		return data, nil
	}

	rv := reflect.New(rt)
	if err := utils.MapToStruct(data, rv.Interface()); err != nil {
		return nil, err
	}
	return rv.Elem().Interface(), nil
}
