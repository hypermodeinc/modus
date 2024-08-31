/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"
	"reflect"

	"hypruntime/utils"
)

func (wa *wasmAdapter) EncodeData(ctx context.Context, typ string, data any) ([]uint64, utils.Cleaner, error) {
	if typ == "" {
		return nil, nil, fmt.Errorf("type is empty")
	}

	return wa.encodeObject(ctx, typ, data)
}

func (wa *wasmAdapter) DecodeData(ctx context.Context, typ string, vals []uint64, pData *any) error {
	if len(vals) == 0 {
		return nil
	}

	if typ == "" {
		return fmt.Errorf("type is empty")
	}

	data, err := wa.decodeObject(ctx, typ, vals)
	if err != nil {
		return err
	}

	// special case for structs represented as maps
	if m, ok := data.(map[string]any); ok {
		if _, ok := (*pData).(map[string]any); !ok {
			return utils.MapToStruct(m, pData)
		}
	}

	// special case for pointers that need to be dereferenced
	if wa.typeInfo.IsPointerType(typ) && reflect.TypeOf(*pData).Kind() != reflect.Ptr {
		*pData = reflect.ValueOf(data).Elem().Interface()
		return nil
	}

	*pData = data
	return nil
}

func (wa *wasmAdapter) GetEncodingLength(ctx context.Context, typ string) (int, error) {

	if wa.typeInfo.IsPrimitiveType(typ) || wa.typeInfo.IsPointerType(typ) || wa.typeInfo.IsMapType(typ) {
		return 1, nil
	}

	if wa.typeInfo.IsStringType(typ) {
		return 2, nil
	}

	if wa.typeInfo.IsSliceType(typ) {
		return 3, nil
	}

	if wa.typeInfo.IsArrayType(typ) {
		return wa.getArrayEncodedLength(ctx, typ)
	}

	return wa.getStructEncodedLength(ctx, typ)
}
