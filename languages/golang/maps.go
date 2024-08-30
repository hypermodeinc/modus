/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"
	"hypruntime/utils"
	"reflect"
)

func (wa *wasmAdapter) encodeMap(ctx context.Context, typ string, obj any) ([]uint64, utils.Cleaner, error) {
	ptr, cln, err := wa.doWriteMap(ctx, typ, obj)
	if err != nil {
		return nil, cln, err
	}

	mapPtr, ok := wa.mod.Memory().ReadUint32Le(ptr)
	if !ok {
		return nil, cln, fmt.Errorf("failed to read map internal pointer from memory")
	}

	return []uint64{uint64(mapPtr)}, cln, nil
}

func (wa *wasmAdapter) decodeMap(ctx context.Context, typ string, vals []uint64) (any, error) {
	mapPtr := uint32(vals[0])
	if mapPtr == 0 {
		return nil, nil
	}

	// we need to make a pointer to the map data
	ptr, cln, err := wa.allocateWasmMemory(ctx, 4)
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

	if ok := wa.mod.Memory().WriteUint32Le(ptr, mapPtr); !ok {
		return nil, fmt.Errorf("failed to write map pointer data to memory")
	}

	// now we can use that pointer to read the map
	return wa.readMap(ctx, typ, ptr)
}

func (wa *wasmAdapter) writeMap(ctx context.Context, typ string, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := wa.doWriteMap(ctx, typ, obj)
	if err != nil {
		return cln, err
	}

	mapPtr, ok := wa.mod.Memory().ReadUint32Le(ptr)
	if !ok {
		return cln, fmt.Errorf("failed to read map internal pointer from memory")
	}

	if ok := wa.mod.Memory().WriteUint32Le(offset, mapPtr); !ok {
		return cln, fmt.Errorf("failed to write map pointer to memory")
	}

	return cln, nil
}

func (wa *wasmAdapter) readMap(ctx context.Context, typ string, offset uint32) (data any, err error) {
	if offset == 0 {
		return nil, nil
	}

	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return nil, err
	}

	res, err := wa.fnReadMap.Call(ctx, uint64(def.Id), uint64(offset))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s from WASM memory: %w", typ, err)
	}
	r := res[0]

	pKeys := uint32(r >> 32)
	pVals := uint32(r)

	keyType, valueType := wa.typeInfo.GetMapSubtypes(typ)
	keys, err := wa.readSlice(ctx, "[]"+keyType, pKeys)
	if err != nil {
		return nil, err
	}
	vals, err := wa.readSlice(ctx, "[]"+valueType, pVals)
	if err != nil {
		return nil, err
	}

	rKeyType, err := wa.getReflectedType(ctx, keyType)
	if err != nil {
		return nil, err
	}
	rValueType, err := wa.getReflectedType(ctx, valueType)
	if err != nil {
		return nil, err
	}

	rvKeys := reflect.ValueOf(keys)
	rvVals := reflect.ValueOf(vals)
	size := rvKeys.Len()

	if rKeyType.Comparable() {
		// return a map
		m := reflect.MakeMapWithSize(reflect.MapOf(rKeyType, rValueType), size)
		for i := 0; i < size; i++ {
			m.SetMapIndex(rvKeys.Index(i), rvVals.Index(i))
		}
		return m.Interface(), nil
	} else {
		// return a pseudo-map
		sliceType := reflect.SliceOf(reflect.StructOf([]reflect.StructField{
			{
				Name: "Key",
				Type: rKeyType,
				Tag:  `json:"key"`,
			},
			{
				Name: "Value",
				Type: rValueType,
				Tag:  `json:"value"`,
			},
		}))
		s := reflect.MakeSlice(sliceType, size, size)
		for i := 0; i < size; i++ {
			s.Index(i).Field(0).Set(rvKeys.Index(i))
			s.Index(i).Field(1).Set(rvVals.Index(i))
		}
		t := reflect.StructOf([]reflect.StructField{
			{
				Name: "Data",
				Type: sliceType,
				Tag:  `json:"$mapdata"`,
			},
		})
		w := reflect.New(t).Elem()
		w.Field(0).Set(s)
		return w.Interface(), nil
	}
}

func (wa *wasmAdapter) doWriteMap(ctx context.Context, typ string, obj any) (pMap uint32, cln utils.Cleaner, err error) {
	if obj == nil {
		return 0, nil, nil
	}

	m, err := utils.ConvertToMap(obj)
	if err != nil {
		return 0, nil, err
	}

	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return 0, nil, err
	}

	pMap, cln, err = wa.makeWasmObject(ctx, def.Id, uint32(len(m)))
	if err != nil {
		return 0, nil, err
	}

	kt, vt := wa.typeInfo.GetMapSubtypes(typ)
	keys, vals := utils.MapKeysAndValues(m)

	cleaner := utils.NewCleanerN(2)
	defer func() {
		// clean up the slices immediately after writing them to the map
		if e := cleaner.Clean(); e != nil && err == nil {
			err = e
		}
	}()

	pKeys, c, err := wa.doWriteSlice(ctx, "[]"+kt, keys)
	cleaner.AddCleaner(c)
	if err != nil {
		return 0, cln, err
	}

	pVals, c, err := wa.doWriteSlice(ctx, "[]"+vt, vals)
	cleaner.AddCleaner(c)
	if err != nil {
		return 0, cln, err
	}

	if _, err := wa.fnWriteMap.Call(ctx, uint64(def.Id), uint64(pMap), uint64(pKeys), uint64(pVals)); err != nil {
		return 0, cln, fmt.Errorf("failed to write %s to WASM memory: %w", typ, err)
	}

	return pMap, cln, nil
}
