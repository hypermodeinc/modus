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

func (p *planner) NewClassHandler(ctx context.Context, typ string, rt reflect.Type) (managedTypeHandler, error) {
	handler := new(classHandler)

	typ = _typeInfo.GetUnderlyingType(typ)
	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	offset := uint32(0)
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
	}
	handler.fieldHandlers = fieldHandlers
	totalSize := offset

	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, totalSize, 0)

	return handler, nil
}

type classHandler struct {
	info          langsupport.TypeHandlerInfo
	typeDef       *metadata.TypeDefinition
	fieldHandlers []langsupport.TypeHandler
}

func (h *classHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *classHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	mapData := make(map[string]any, len(h.fieldHandlers))
	fieldOffset := uint32(0)
	for i, field := range h.typeDef.Fields {
		handler := h.fieldHandlers[i]
		size := handler.Info().TypeSize()
		fieldOffset += langsupport.GetAlignmentPadding(fieldOffset, size)

		val, err := handler.Read(ctx, wa, offset+fieldOffset)
		if err != nil {
			return nil, err
		}

		mapData[field.Name] = val
		fieldOffset += size
	}

	rtData := h.info.RuntimeType()
	if rtData.Kind() == reflect.Map {
		return mapData, nil
	}

	rvDataPtr := reflect.New(rtData)
	if err := utils.MapToStruct(mapData, rvDataPtr.Interface()); err != nil {
		return nil, err
	}

	return rvDataPtr.Elem().Interface(), nil
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
		c, err := handler.Write(ctx, wa, offset+fieldOffset, fieldObj)
		cln.AddCleaner(c)
		if err != nil {
			return cln, err
		}

		fieldOffset += size
	}

	return cln, nil
}
