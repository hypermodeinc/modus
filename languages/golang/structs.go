/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"hypruntime/logger"
	"hypruntime/utils"
)

const maxDepth = 5 // TODO: make this based on the depth requested in the query

func (wa *wasmAdapter) encodeStruct(ctx context.Context, typ string, obj any) ([]uint64, utils.Cleaner, error) {

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

	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return nil, nil, err
	}

	numFields := len(def.Fields)
	results := make([]uint64, 0, numFields*2)
	cleaner := utils.NewCleanerN(numFields)

	for _, field := range def.Fields {
		var fieldObj any
		if mapObj != nil {
			// case sensitive when reading from map
			fieldObj = mapObj[field.Name]
		} else {
			// case insensitive when reading from struct
			fieldObj = rvObj.FieldByNameFunc(func(s string) bool { return strings.EqualFold(s, field.Name) }).Interface()
		}

		vals, cln, err := wa.encodeObject(ctx, field.Type, fieldObj)
		cleaner.AddCleaner(cln)
		if err != nil {
			return nil, cleaner, err
		}
		results = append(results, vals...)
	}

	return results, cleaner, nil
}

func (wa *wasmAdapter) decodeStruct(ctx context.Context, typ string, vals []uint64) (any, error) {
	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return nil, err
	}
	if len(def.Fields) == 1 {
		data, err := wa.decodeObject(ctx, def.Fields[0].Type, vals)
		if err != nil {
			return nil, err
		}

		m := map[string]any{def.Fields[0].Name: data}
		return getStructOutput(ctx, wa, typ, m)
	}

	// TODO: Implement decoding of structs with multiple fields from multiple values
	return nil, fmt.Errorf("decoding struct of type %s is not supported", typ)
}

func (wa *wasmAdapter) getStructEncodedLength(ctx context.Context, typ string) (int, error) {
	total := 0

	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return 0, err
	}

	for _, field := range def.Fields {
		i, err := wa.GetEncodingLength(ctx, field.Type)
		if err != nil {
			return 0, err
		}
		total += i
	}

	return total, nil
}

func (wa *wasmAdapter) writeStruct(ctx context.Context, typ string, offset uint32, obj any) (utils.Cleaner, error) {
	size, err := wa.typeInfo.GetSizeOfType(ctx, typ)
	if err != nil {
		return nil, err
	}

	ptr, cln, err := wa.doWriteStruct(ctx, typ, obj)
	if err != nil {
		return cln, err
	}

	// copy the struct data to the target memory location
	mem := wa.mod.Memory()
	if bytes, ok := mem.Read(ptr, size); !ok {
		return cln, errors.New("failed to read struct data")
	} else if !mem.Write(offset, bytes) {
		return cln, errors.New("failed to write struct data")
	}

	return cln, nil
}

func (wa *wasmAdapter) readStruct(ctx context.Context, typ string, offset uint32) (any, error) {

	// Check for recursion
	if wa.visitedPtrs[offset] >= maxDepth {
		logger.Warn(ctx).Bool("user_visible", true).Msgf("Excessive recursion detected in %s. Stopping at depth %d.", typ, maxDepth)
		return nil, nil
	}
	wa.visitedPtrs[offset]++
	defer func() {
		n := wa.visitedPtrs[offset]
		if n == 1 {
			delete(wa.visitedPtrs, offset)
		} else {
			wa.visitedPtrs[offset] = n - 1
		}
	}()

	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return nil, err
	}

	data := make(map[string]any, len(def.Fields))
	fieldOffset := uint32(0)
	for _, field := range def.Fields {
		size, err := wa.typeInfo.GetSizeOfType(ctx, field.Type)
		if err != nil {
			return nil, err
		}

		fieldOffset += wa.typeInfo.getAlignmentPadding(fieldOffset, size)

		val, err := wa.readObject(ctx, field.Type, offset+fieldOffset)
		if err != nil {
			return nil, err
		}
		data[field.Name] = val

		fieldOffset += size
	}

	return getStructOutput(ctx, wa, typ, data)
}

func (wa *wasmAdapter) doWriteStruct(ctx context.Context, typ string, obj any) (uint32, utils.Cleaner, error) {

	var mapObj map[string]any
	var rvObj reflect.Value
	if m, ok := obj.(map[string]any); ok {
		mapObj = m
	} else {
		rvObj = reflect.ValueOf(obj)
		if rvObj.Kind() == reflect.Ptr {
			rvObj = rvObj.Elem()
		}
		if rvObj.Kind() != reflect.Struct {
			return 0, nil, fmt.Errorf("expected a struct, got %s", rvObj.Kind())
		}
	}

	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return 0, nil, err
	}

	ptr, cln, err := wa.newWasmObject(ctx, def.Id)
	if err != nil {
		return 0, cln, err
	}

	cleaner := utils.NewCleanerN(len(def.Fields) + 1)
	cleaner.AddCleaner(cln)

	fieldOffset := uint32(0)
	for _, field := range def.Fields {
		size, err := wa.typeInfo.GetSizeOfType(ctx, field.Type)
		if err != nil {
			return 0, cleaner, err
		}

		fieldOffset += wa.typeInfo.getAlignmentPadding(fieldOffset, size)

		var val any
		if mapObj != nil {
			// case sensitive when reading from map
			val = mapObj[field.Name]
		} else {
			// case insensitive when reading from struct
			val = rvObj.FieldByNameFunc(func(s string) bool { return strings.EqualFold(s, field.Name) }).Interface()
		}

		cln, err := wa.writeObject(ctx, field.Type, ptr+fieldOffset, val)
		cleaner.AddCleaner(cln)
		if err != nil {
			return 0, cleaner, err
		}

		fieldOffset += size
	}

	return ptr, cleaner, nil
}

func getStructOutput(ctx context.Context, wa *wasmAdapter, typ string, data map[string]any) (any, error) {
	rt, err := wa.getReflectedType(ctx, typ)
	if err != nil {
		return nil, err
	}
	if rt.Kind() == reflect.Struct {
		rv := reflect.New(rt)
		if err := utils.MapToStruct(data, rv.Interface()); err != nil {
			return nil, err
		}
		return rv.Elem().Interface(), nil
	}

	return data, nil
}
