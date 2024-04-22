/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func EncodeValue(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val any) (uint64, error) {
	switch typ.Path {
	case "bool":
		b, ok := val.(bool)
		if !ok {
			return 0, fmt.Errorf("input value is not a bool")
		}

		if b {
			return 1, nil
		} else {
			return 0, nil
		}

	case "i8", "i16", "i32", "u8", "u16", "u32":
		var x int32
		switch v := val.(type) {
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
		case uint:
			x = int32(v)
		case uint8:
			x = int32(v)
		case uint16:
			x = int32(v)
		case uint32:
			x = int32(v)
		default:
			return 0, fmt.Errorf("input value is not an int")
		}

		return wasm.EncodeI32(x), nil

	case "i64", "u64":
		var x int64
		switch v := val.(type) {
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
		switch v := val.(type) {
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
		switch v := val.(type) {
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
	offset, err := writeObject(ctx, mod, typ, val)
	return uint64(offset), err
}

func DecodeValueAs[T any](ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val uint64) (T, error) {
	var result T
	r, err := DecodeValue(ctx, mod, typ, val)
	if err != nil {
		return result, err
	}

	switch v := r.(type) {
	case T:
		return v, nil
	case map[string]any:
		out, err := utils.ConvertToStruct[T](v)
		if err != nil {
			return result, err
		}
		return out, nil
	default:
		return result, fmt.Errorf("unexpected type %T, expected %T", r, result)
	}
}

func DecodeValue(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val uint64) (any, error) {

	// Handle null values if the type is nullable
	if isNullable(typ.Path) {
		if val == 0 {
			return nil, nil
		}
		typ = plugins.TypeInfo{
			Name: typ.Name[:len(typ.Name)-7], // remove " | null"
			Path: typ.Path[:len(typ.Path)-5], // remove "|null"
		}
	}

	// Primitive types
	switch typ.Path {
	case "void":
		return 0, nil

	case "bool":
		return val != 0, nil

	case "i8", "i16", "i32", "u8", "u16", "u32":
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
