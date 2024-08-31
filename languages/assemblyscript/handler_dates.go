/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

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

// Reference: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/date.ts

func (p *planner) NewDateHandler(typ string, rt reflect.Type) (managedTypeHandler, error) {
	handler := new(dateHandler)
	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, 20, 0)

	typ = _typeInfo.GetUnderlyingType(typ)
	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	return handler, nil
}

type dateHandler struct {
	info    langsupport.TypeHandlerInfo
	typeDef *metadata.TypeDefinition
}

func (h *dateHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *dateHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	val, ok := wa.Memory().ReadUint64Le(offset + 16)
	if !ok {
		return nil, errors.New("error reading timestamp from wasm memory")
	}

	ts := int64(val)
	tm := time.UnixMilli(ts).UTC()
	return tm, nil
}

func (h *dateHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	var tm time.Time
	switch t := obj.(type) {
	case time.Time:
		tm = t.UTC()
	case *time.Time:
		tm = t.UTC()
	case utils.JSONTime:
		tm = time.Time(t).UTC()
	case *utils.JSONTime:
		tm = time.Time(*t).UTC()
	default:
		return nil, fmt.Errorf("incompatible type for Date object: %T", obj)
	}

	if ok := wa.Memory().WriteUint32Le(offset, uint32(tm.Year())); !ok {
		return nil, errors.New("failed to write Date object's year to WASM memory")
	}

	if ok := wa.Memory().WriteUint32Le(offset+4, uint32(tm.Month())); !ok {
		return nil, errors.New("failed to write Date object's month to WASM memory")
	}

	if ok := wa.Memory().WriteUint32Le(offset+8, uint32(tm.Day())); !ok {
		return nil, errors.New("failed to write Date object's day to WASM memory")
	}

	if ok := wa.Memory().WriteUint64Le(offset+16, uint64(tm.UnixMilli())); !ok {
		return nil, errors.New("failed to write Date object's timestamp to WASM memory")
	}

	return nil, nil
}
