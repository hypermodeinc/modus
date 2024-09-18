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

func (p *planner) NewMapHandler(ctx context.Context, typ string, rt reflect.Type) (langsupport.TypeHandler, error) {
	handler := NewTypeHandler[mapHandler](p, typ)

	// maps are passed by reference using a 4 byte pointer
	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)

	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	keyType, valueType := _typeInfo.GetMapSubtypes(typ)
	keysHandler, err := p.GetHandler(ctx, "[]"+keyType)
	if err != nil {
		return nil, err
	}
	handler.keysHandler = keysHandler

	valuesHandler, err := p.GetHandler(ctx, "[]"+valueType)
	if err != nil {
		return nil, err
	}
	handler.valuesHandler = valuesHandler

	rtKey := keysHandler.Info().RuntimeType().Elem()
	rtValue := valuesHandler.Info().RuntimeType().Elem()
	if !rtKey.Comparable() {
		handler.usePseudoMap = true
		handler.rtPseudoMapSlice = reflect.SliceOf(reflect.StructOf([]reflect.StructField{
			{
				Name: "Key",
				Type: rtKey,
				Tag:  `json:"key"`,
			},
			{
				Name: "Value",
				Type: rtValue,
				Tag:  `json:"value"`,
			},
		}))

		handler.rtPseudoMap = reflect.StructOf([]reflect.StructField{
			{
				Name: "Data",
				Type: handler.rtPseudoMapSlice,
				Tag:  `json:"$mapdata"`,
			},
		})
	}

	return handler, nil
}

type mapHandler struct {
	info             langsupport.TypeHandlerInfo
	typeDef          *metadata.TypeDefinition
	keysHandler      langsupport.TypeHandler
	valuesHandler    langsupport.TypeHandler
	usePseudoMap     bool
	rtPseudoMap      reflect.Type
	rtPseudoMapSlice reflect.Type
}

func (h *mapHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *mapHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	res, err := wa.(*wasmAdapter).fnReadMap.Call(ctx, uint64(h.typeDef.Id), uint64(offset))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s from WASM memory: %w", h.Info().TypeName(), err)
	}
	r := res[0]

	pKeys := uint32(r >> 32)
	pVals := uint32(r)

	keys, err := h.keysHandler.Read(ctx, wa, pKeys)
	if err != nil {
		return nil, err
	}
	vals, err := h.valuesHandler.Read(ctx, wa, pVals)
	if err != nil {
		return nil, err
	}

	rvKeys := reflect.ValueOf(keys)
	rvVals := reflect.ValueOf(vals)
	size := rvKeys.Len()

	rtKey := h.keysHandler.Info().RuntimeType().Elem()
	if rtKey.Comparable() {
		// return a map
		m := reflect.MakeMapWithSize(h.info.RuntimeType(), size)
		for i := 0; i < size; i++ {
			m.SetMapIndex(rvKeys.Index(i), rvVals.Index(i))
		}
		return m.Interface(), nil
	} else {
		s := reflect.MakeSlice(h.rtPseudoMapSlice, size, size)
		for i := 0; i < size; i++ {
			s.Index(i).Field(0).Set(rvKeys.Index(i))
			s.Index(i).Field(1).Set(rvVals.Index(i))
		}

		m := reflect.New(h.rtPseudoMap).Elem()
		m.Field(0).Set(s)
		return m.Interface(), nil
	}
}

func (h *mapHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := h.doWriteMap(ctx, wa, obj)
	if err != nil {
		return cln, err
	}

	mapPtr, ok := wa.Memory().ReadUint32Le(ptr)
	if !ok {
		return cln, errors.New("failed to read map internal pointer from memory")
	}

	if ok := wa.Memory().WriteUint32Le(offset, mapPtr); !ok {
		return cln, errors.New("failed to write map pointer to memory")
	}

	return cln, nil
}

func (h *mapHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	mapPtr := uint32(vals[0])
	if mapPtr == 0 {
		return nil, nil
	}

	// we need to make a pointer to the map data
	ptr, cln, err := wa.(*wasmAdapter).AllocateMemory(ctx, 4)
	defer func() {
		if cln != nil {
			if e := cln.Clean(); e != nil && err == nil {
				err = e
			}
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for map pointer: %w", err)
	}

	if ok := wa.Memory().WriteUint32Le(ptr, mapPtr); !ok {
		return nil, errors.New("failed to write map pointer data to memory")
	}

	// now we can use that pointer to read the map
	return h.Read(ctx, wa, ptr)
}

func (h *mapHandler) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	ptr, cln, err := h.doWriteMap(ctx, wa, obj)
	if err != nil {
		return nil, cln, err
	}

	mapPtr, ok := wa.Memory().ReadUint32Le(ptr)
	if !ok {
		return nil, cln, errors.New("failed to read map internal pointer from memory")
	}

	return []uint64{uint64(mapPtr)}, cln, nil
}

func (h *mapHandler) doWriteMap(ctx context.Context, wa langsupport.WasmAdapter, obj any) (pMap uint32, cln utils.Cleaner, err error) {
	if utils.HasNil(obj) {
		return 0, nil, nil
	}

	m, err := utils.ConvertToMap(obj)
	if err != nil {
		return 0, nil, err
	}
	keys, vals := utils.MapKeysAndValues(m)

	pMap, cln, err = wa.(*wasmAdapter).makeWasmObject(ctx, h.typeDef.Id, uint32(len(m)))
	if err != nil {
		return 0, nil, err
	}

	innerCln := utils.NewCleanerN(2)
	defer func() {
		// clean up the slices after the map is written to memory
		if e := innerCln.Clean(); e != nil && err == nil {
			err = e
		}
	}()

	pKeys, c, err := h.keysHandler.(sliceWriter).doWriteSlice(ctx, wa, keys)
	innerCln.AddCleaner(c)
	if err != nil {
		return 0, cln, err
	}

	pVals, c, err := h.valuesHandler.(sliceWriter).doWriteSlice(ctx, wa, vals)
	innerCln.AddCleaner(c)
	if err != nil {
		return 0, cln, err
	}

	if _, err := wa.(*wasmAdapter).fnWriteMap.Call(ctx, uint64(h.typeDef.Id), uint64(pMap), uint64(pKeys), uint64(pVals)); err != nil {
		return 0, cln, fmt.Errorf("failed to write %s to WASM memory: %w", h.info.TypeName(), err)
	}

	return pMap, cln, nil
}
