/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func EncodeValue(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val any) (uint64, error) {

	switch typ.Name {

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
		n, err := val.(json.Number).Int64()
		if err != nil {
			return 0, fmt.Errorf("input value is not an int")
		}

		return wasm.EncodeI32(int32(n)), nil

	case "i64", "u64":
		n, err := val.(json.Number).Int64()
		if err != nil {
			return 0, fmt.Errorf("input value is not an int")
		}

		return wasm.EncodeI64(n), nil

	case "f32":
		n, err := val.(json.Number).Float64()
		if err != nil {
			return 0, fmt.Errorf("input value is not a float")
		}

		return wasm.EncodeF32(float32(n)), nil

	case "f64":
		n, err := val.(json.Number).Float64()
		if err != nil {
			return 0, fmt.Errorf("input value is not a float")
		}

		return wasm.EncodeF64(n), nil

	case "string":
		s, ok := val.(string)
		if !ok {
			return 0, fmt.Errorf("input value is not a string")
		}

		// Note, strings are passed as a pointer to a string in wasm memory
		ptr := WriteString(ctx, mod, s)
		return uint64(ptr), nil

	case "Date":
		var t time.Time
		switch v := val.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, err
			}
			t = time.UnixMilli(n)
		case string:
			var err error
			t, err = utils.ParseTime(v)
			if err != nil {
				return 0, err
			}
		}

		ptr, err := writeDate(ctx, mod, t)
		if err != nil {
			return 0, err
		}

		return uint64(ptr), nil

	// TODO: custom types

	default:
		return 0, fmt.Errorf("unknown parameter type: %s", typ)
	}
}

func DecodeValue(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val uint64) (any, error) {

	// Primitive types
	switch typ.Name {
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

	// Managed types, held in wasm memory
	mem := mod.Memory()
	return readObject(ctx, mem, typ, uint32(val))
}

func getTypeDefinition(ctx context.Context, typePath string) (plugins.TypeDefinition, error) {
	plugin := ctx.Value(utils.PluginContextKey).(*plugins.Plugin)

	types := plugin.Types
	info, ok := types[typePath]
	if !ok {
		return plugins.TypeDefinition{}, fmt.Errorf("info for type %s not found in plugin %s", typePath, plugin.Name())
	}

	return info, nil
}
