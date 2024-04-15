/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hmruntime/host"
	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func CallFunction(ctx context.Context, mod wasm.Module, info FunctionInfo, inputs map[string]any) (any, error) {
	fnName := info.Function.Name
	logger.Info(ctx).
		Str("function", fnName).
		Bool("user_visible", true).
		Msg("Calling function.")

	start := time.Now()
	result, err := callFunction(ctx, mod, info, inputs)
	duration := time.Since(start)

	if err != nil {
		logger.Err(ctx, err).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Error while executing function.")
	} else {
		logger.Info(ctx).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Function completed successfully.")
	}

	return result, err
}

func callFunction(ctx context.Context, mod wasm.Module, info FunctionInfo, inputs map[string]any) (any, error) {
	fnName := info.Function.Name

	// Get the function
	fn := mod.ExportedFunction(fnName)
	if fn == nil {
		return nil, fmt.Errorf("function %s not found in plugin %s", fnName, info.Plugin.Name())
	}

	// Get parameters to pass as input to the plugin function
	params := make([]uint64, len(info.Function.Parameters))
	for i, arg := range info.Function.Parameters {
		val := inputs[arg.Name]
		if val == nil {
			return nil, fmt.Errorf("parameter '%s' is missing", arg.Name)
		}

		param, err := convertParam(ctx, mod, arg.Type, val)
		if err != nil {
			return nil, fmt.Errorf("function parameter '%s' is invalid: %w", arg.Name, err)
		}

		params[i] = param
	}

	// Call the wasm function
	res, err := fn.Call(ctx, params...)
	if err != nil {
		return nil, err
	}

	// Get the result
	result, err := convertResult(ctx, mod, info.Function.ReturnType, res[0])
	if err != nil {
		return nil, fmt.Errorf("function result is invalid: %w", err)
	}

	return result, nil
}

func convertParam(ctx context.Context, mod wasm.Module, asType plugins.TypeInfo, val any) (uint64, error) {

	switch asType.Name {

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
		ptr := writeString(ctx, mod, s)
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
		return 0, fmt.Errorf("unknown parameter type: %s", asType)
	}
}

func convertResult(ctx context.Context, mod wasm.Module, asType plugins.TypeInfo, res uint64) (any, error) {

	// Primitive types
	switch asType.Name {
	case "void":
		return 0, nil

	case "bool":
		return res != 0, nil

	case "i8", "i16", "i32", "u8", "u16", "u32":
		return wasm.DecodeI32(res), nil

	case "i64", "u64":
		return int64(res), nil

	case "f32":
		return wasm.DecodeF32(res), nil

	case "f64":
		return wasm.DecodeF64(res), nil
	}

	// Handled managed types
	mem := mod.Memory()
	switch asType.Name {
	case "string":
		return readString(mem, uint32(res))

	case "Date":
		return readDate(mem, uint32(res))
	}

	// Custom managed types
	return readObject(ctx, mem, asType, uint32(res))
}

func getTypeDefinition(ctx context.Context, typePath string) (plugins.TypeDefinition, error) {
	pluginName := ctx.Value(utils.PluginNameContextKey).(string)
	plugin, ok := host.Plugins.GetByName(pluginName)
	if !ok {
		return plugins.TypeDefinition{}, fmt.Errorf("plugin %s not found", pluginName)
	}

	types := plugin.Types
	info, ok := types[typePath]
	if !ok {
		return plugins.TypeDefinition{}, fmt.Errorf("info for type %s not found in plugin %s", typePath, pluginName)
	}

	return info, nil
}
