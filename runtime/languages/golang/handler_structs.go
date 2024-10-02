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

func (p *planner) NewStructHandler(ctx context.Context, ti langsupport.TypeInfo) (langsupport.TypeHandler, error) {
	handler := &structHandler{
		typeHandler: *NewTypeHandler(ti),
	}
	p.AddHandler(handler)

	typeDef, err := p.metadata.GetTypeDefinition(ti.Name())
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	fieldTypes := ti.ObjectFieldTypes()
	fieldHandlers := make([]langsupport.TypeHandler, len(fieldTypes))
	for i, fieldType := range fieldTypes {
		fieldHandler, err := p.GetHandler(ctx, fieldType.Name())
		if err != nil {
			return nil, err
		}
		fieldHandlers[i] = fieldHandler
	}

	handler.fieldHandlers = fieldHandlers
	return handler, nil
}

type structHandler struct {
	typeHandler
	typeDef       *metadata.TypeDefinition
	fieldHandlers []langsupport.TypeHandler
}

func (h *structHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {

	// Check for recursion
	visitedPtrs := wa.(*wasmAdapter).visitedPtrs
	if visitedPtrs[offset] >= maxDepth {
		logger.Warn(ctx).Bool("user_visible", true).Msgf("Excessive recursion detected in %s. Stopping at depth %d.", h.typeInfo.Name(), maxDepth)
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

	fieldOffsets := h.typeInfo.ObjectFieldOffsets()

	m := make(map[string]any, len(h.fieldHandlers))
	for i, field := range h.typeDef.Fields {
		handler := h.fieldHandlers[i]
		fieldOffset := offset + fieldOffsets[i]
		val, err := handler.Read(ctx, wa, fieldOffset)
		if err != nil {
			return nil, err
		}
		m[field.Name] = val
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
	fieldOffsets := h.typeInfo.ObjectFieldOffsets()
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

		fieldOffset := offset + fieldOffsets[i]
		handler := h.fieldHandlers[i]
		cln, err := handler.Write(ctx, wa, fieldOffset, fieldObj)
		cleaner.AddCleaner(cln)
		if err != nil {
			return cleaner, err
		}
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
	return nil, fmt.Errorf("decoding struct of type %s is not supported", h.typeInfo.Name())
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

func (h *structHandler) getStructOutput(data map[string]any) (any, error) {
	rt := h.typeInfo.ReflectedType()
	if rt.Kind() == reflect.Map {
		return data, nil
	}

	rv := reflect.New(rt)
	if err := utils.MapToStruct(data, rv.Interface()); err != nil {
		return nil, err
	}
	return rv.Elem().Interface(), nil
}
