/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"
	"reflect"

	"hypruntime/utils"

	"github.com/spf13/cast"
)

func (wa *wasmAdapter) encodePointer(ctx context.Context, typ string, obj any) ([]uint64, utils.Cleaner, error) {
	ptr, cln, err := wa.doWritePointer(ctx, typ, obj)
	if err != nil {
		return nil, cln, err
	}
	return []uint64{uint64(ptr)}, cln, nil
}

func (wa *wasmAdapter) decodePointer(ctx context.Context, typ string, vals []uint64) (any, error) {
	dataType := typ[1:]
	data, err := wa.readObject(ctx, dataType, uint32(vals[0]))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	// special case for structs represented as maps
	rt := reflect.TypeOf(data)
	if rt.Kind() == reflect.Map && !wa.typeInfo.IsMapType(dataType) {
		return data, nil
	}

	ptr := reflect.New(rt)
	ptr.Elem().Set(reflect.ValueOf(data))
	return ptr.Interface(), nil
}

func (wa *wasmAdapter) writePointer(ctx context.Context, typ string, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := wa.doWritePointer(ctx, typ, obj)
	if err != nil {
		return cln, err
	}
	if ok := wa.mod.Memory().WriteUint32Le(offset, ptr); !ok {
		return cln, fmt.Errorf("failed to write object pointer to memory")
	}
	return cln, nil
}

func (wa *wasmAdapter) readPointer(ctx context.Context, typ string, offset uint32) (any, error) {
	ptr, ok := wa.mod.Memory().ReadUint32Le(offset)
	if !ok {
		return nil, fmt.Errorf("failed to read object pointer from memory")
	}

	t := wa.typeInfo.GetUnderlyingType(typ)
	data, err := wa.readObject(ctx, t, ptr)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(data)
	p := reflect.New(rv.Type())
	p.Elem().Set(rv)
	return p.Interface(), nil
}

func (wa *wasmAdapter) doWritePointer(ctx context.Context, typ string, obj any) (uint32, utils.Cleaner, error) {
	if obj == nil {
		return 0, nil, nil
	}

	// dereference the pointer (if necessary)
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return 0, nil, nil
		} else {
			obj = rv.Elem().Interface()
		}
	}

	// get the underlying type
	t := wa.typeInfo.GetUnderlyingType(typ)

	// handle pointers to strings
	if wa.typeInfo.IsStringType(t) {
		str, err := cast.ToStringE(obj)
		if err != nil {
			return 0, nil, err
		}
		return wa.doWriteString(ctx, str)
	}

	// handle pointers to slices
	if wa.typeInfo.IsSliceType(t) {
		return wa.doWriteSlice(ctx, t, obj)
	}

	// handle pointers to maps
	if wa.typeInfo.IsMapType(t) {
		return wa.doWriteMap(ctx, t, obj)
	}

	// handle pointers to primitive types
	if wa.typeInfo.IsPrimitiveType(t) {
		def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
		if err != nil {
			return 0, nil, err
		}

		ptr, cln, err := wa.newWasmObject(ctx, def.Id)
		if err != nil {
			return 0, cln, nil
		}

		c, err := wa.writeObject(ctx, t, ptr, obj)
		cln.AddCleaner(c)
		if err != nil {
			return 0, cln, err
		}

		return ptr, cln, nil
	}

	// handle pointer to time.Time types
	if wa.typeInfo.IsTimestampType(t) {
		return wa.doWriteTime(ctx, obj)
	}

	// handle pointers to structs
	return wa.doWriteStruct(ctx, t, obj)
}
