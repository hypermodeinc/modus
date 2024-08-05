/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hmruntime/plugins/metadata"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	wasm "github.com/tetratelabs/wazero/api"
)

func EncodeValue[T any](ctx context.Context, mod wasm.Module, data T) (uint64, error) {
	typ, err := getTypeInfoForType[T]()
	if err != nil {
		return 0, err
	}

	return encodeValue(ctx, mod, typ, data)
}

func encodeValue(ctx context.Context, mod wasm.Module, typ metadata.TypeInfo, data any) (uint64, error) {
	// For most calls, we don't need to pin the memory.
	// If it needs to be pinned, the caller will do it.
	return doEncodeValue(ctx, mod, typ, data, false)
}

func EncodeValueForParameter(ctx context.Context, mod wasm.Module, typ metadata.TypeInfo, data any) (uint64, error) {
	// For the inbound parameters, we need to pin the memory.
	// Otherwise, just allocating more parameters could cause the GC to run and free the memory we just allocated.
	// Note we don't bother tracking these to unpin later, because we discard the module instance and its memory after the call.
	return doEncodeValue(ctx, mod, typ, data, true)
}

func doEncodeValue(ctx context.Context, mod wasm.Module, typ metadata.TypeInfo, data any, pin bool) (val uint64, err error) {

	// Recover from panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			}
		}
	}()

	// Handle null values if the type is nullable
	if isNullable(typ.Path) {
		if data == nil {
			return 0, nil
		}
		rv := reflect.ValueOf(data)
		switch rv.Kind() {
		case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map:
			if rv.IsNil() {
				return 0, nil
			}
		}

		typ = removeNull(typ)
	}

	// Primitive types
	switch typ.Path {
	case "bool":
		b, ok := data.(bool)
		if !ok {
			return 0, fmt.Errorf("input value is not a bool")
		}

		if b {
			return 1, nil
		} else {
			return 0, nil
		}

	case "u8", "u16", "u32":
		var x uint32
		switch v := data.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, fmt.Errorf("input value is not an int")
			}
			x = uint32(n)
		case uint:
			x = uint32(v)
		case uint8:
			x = uint32(v)
		case uint16:
			x = uint32(v)
		case uint32:
			x = uint32(v)
		default:
			return 0, fmt.Errorf("input value is not an unsigned integer")
		}

		return wasm.EncodeU32(x), nil

	case "i8", "i16", "i32":
		var x int32
		switch v := data.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, fmt.Errorf("input value is not an int")
			}
			x = int32(n)
		case int:
			x = int32(v)
		case int8:
			x = int32(v)
		case int16:
			x = int32(v)
		case int32:
			x = v
		default:
			return 0, fmt.Errorf("input value is not an signed integer")
		}

		return wasm.EncodeI32(x), nil

	case "i64", "u64":
		var x int64
		switch v := data.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, fmt.Errorf("input value is not an int")
			}
			x = int64(n)
		case int64:
			x = v
		case uint64:
			x = int64(v)
		default:
			return 0, fmt.Errorf("input value is not an int")
		}

		return wasm.EncodeI64(x), nil

	case "f32":
		var x float32
		switch v := data.(type) {
		case json.Number:
			n, err := v.Float64()
			if err != nil {
				return 0, fmt.Errorf("input value is not a float")
			}
			x = float32(n)
		case float32:
			x = v
		default:
			return 0, fmt.Errorf("input value is not a float")
		}

		return wasm.EncodeF32(x), nil

	case "f64":
		var x float64
		switch v := data.(type) {
		case json.Number:
			n, err := v.Float64()
			if err != nil {
				return 0, fmt.Errorf("input value is not a float")
			}
			x = float64(n)
		case float64:
			x = v
		default:
			return 0, fmt.Errorf("input value is not a float")
		}

		return wasm.EncodeF64(x), nil
	}

	// Managed types need to be written to wasm memory
	offset, err := writeObject(ctx, mod, typ, data)
	if err != nil {
		return 0, err
	}

	// Pin the memory if requested
	if pin {
		err := pinWasmMemory(ctx, mod, offset)
		if err != nil {
			return 0, fmt.Errorf("failed to pin wasm memory: %w", err)
		}
	}

	return uint64(offset), nil
}

func DecodeValueAs[T any](ctx context.Context, mod wasm.Module, val uint64) (data T, err error) {

	// Recover from panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			}
		}
	}()

	var result T

	typ, err := getTypeInfoForType[T]()
	if err != nil {
		return result, err
	}

	r, err := DecodeValue(ctx, mod, typ, val)
	if err != nil {
		return result, err
	}

	switch v := r.(type) {
	case T:
		// If the type is already the expected type, return it.
		return v, nil
	case map[string]any:
		// If the type is a map, convert it to the expected type.
		// This is expected in the case of a host function that takes a struct as a parameter.
		return mapToStruct[T](v)
	}

	return result, fmt.Errorf("unexpected type %T, expected %T", r, result)
}

func mapToStruct[T any](m map[string]any) (T, error) {
	var result T

	config := &mapstructure.DecoderConfig{
		Result: &result,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return result, err
	}

	err = decoder.Decode(m)
	return result, err
}

func DecodeValue(ctx context.Context, mod wasm.Module, typ metadata.TypeInfo, val uint64) (data any, err error) {

	// Recover from panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			}
		}
	}()

	// Handle null values if the type is nullable
	if isNullable(typ.Path) {
		if val == 0 {
			return nil, nil
		}
		typ = removeNull(typ)
	}

	// Primitive types
	switch typ.Path {
	case "void":
		return 0, nil

	case "bool":
		return val != 0, nil

	case "u8", "u16", "u32":
		return wasm.DecodeU32(val), nil

	case "i8", "i16", "i32":
		return wasm.DecodeI32(val), nil

	case "i64", "u64":
		return int64(val), nil

	case "f32":
		return wasm.DecodeF32(val), nil

	case "f64":
		return wasm.DecodeF64(val), nil
	}

	// Managed types are read from wasm memory
	mem := mod.Memory()
	return readObject(ctx, mem, typ, uint32(val))
}

func removeNull(typ metadata.TypeInfo) metadata.TypeInfo {
	return metadata.TypeInfo{
		Name: typ.Name[:len(typ.Name)-7], // remove " | null"
		Path: typ.Path[:len(typ.Path)-5], // remove "|null"
	}
}
