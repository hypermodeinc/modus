/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"hypruntime/langsupport"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

func (p *planner) NewClassHandler(ctx context.Context, ti langsupport.TypeInfo) (managedTypeHandler, error) {

	handler := &classHandler{
		typeHandler: *NewTypeHandler(ti),
	}

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

type classHandler struct {
	typeHandler
	typeDef       *metadata.TypeDefinition
	fieldHandlers []langsupport.TypeHandler
}

func (h *classHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

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

func (h *classHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	var mapObj map[string]any
	var rvObj reflect.Value
	if m, ok := obj.(map[string]any); ok {
		mapObj = m
	} else {
		rvObj = reflect.ValueOf(obj)
		if rvObj.Kind() == reflect.Pointer {
			rvObj = rvObj.Elem()
		}
		if rvObj.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected a struct, got %s", rvObj.Kind())
		}
	}

	cln := utils.NewCleanerN(len(h.fieldHandlers))

	fieldOffsets := h.typeInfo.ObjectFieldOffsets()
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
		c, err := handler.Write(ctx, wa, fieldOffset, fieldObj)
		cln.AddCleaner(c)
		if err != nil {
			return cln, err
		}
	}

	return cln, nil
}

func (h *classHandler) getStructOutput(data map[string]any) (any, error) {
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
