/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hmruntime/logger"
	"hmruntime/schema"

	"github.com/dgraph-io/gqlparser/ast"
	wasm "github.com/tetratelabs/wazero/api"
)

func CallFunction(ctx context.Context, mod wasm.Module, info schema.FunctionInfo, inputs map[string]any) (any, error) {
	fnName := info.FunctionName()
	fn := mod.ExportedFunction(fnName)
	def := fn.Definition()
	paramTypes := def.ParamTypes()
	schema := info.Schema

	// Get parameters to pass as input to the plugin function
	// Note that we can't use def.ParamNames() because they are only available when the plugin
	// is compiled in debug mode. They're striped by optimization in release mode.
	// Instead, we can use the argument names from the schema.
	// Also note, that the order of the arguments from schema should match order of params in wasm.
	params := make([]uint64, len(paramTypes))
	for i, arg := range schema.FunctionArgs(ctx) {
		val := inputs[arg.Name]
		if val == nil {
			return nil, fmt.Errorf("parameter %s is missing", arg.Name)
		}

		param, err := convertParam(ctx, mod, *arg.Type, paramTypes[i], val)
		if err != nil {
			return nil, fmt.Errorf("parameter %s is invalid: %w", arg.Name, err)
		}

		params[i] = param
	}

	// Call the wasm function
	logger.Info(ctx).
		Str("plugin", info.PluginName).
		Str("function", fnName).
		Str("resolver", schema.Resolver()).
		Bool("user_visible", true).
		Msg("Calling function.")
	start := time.Now()
	res, err := fn.Call(ctx, params...)
	duration := time.Since(start)
	if err != nil {
		logger.Err(ctx, err).
			Str("plugin", info.PluginName).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Error while executing function.")

		return nil, err
	}

	// Get the result
	mem := mod.Memory()
	result, err := convertResult(mem, *schema.FieldDef.Type, def.ResultTypes()[0], res[0])
	if err != nil {
		logger.Err(ctx, err).
			Str("plugin", info.PluginName).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Str("schema_type", schema.FieldDef.Type.NamedType).
			Bool("user_visible", true).
			Msg("Failed to convert the result of the function to the schema field type.")

		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	logger.Info(ctx).
		Str("plugin", info.PluginName).
		Str("function", fnName).
		Dur("duration_ms", duration).
		Bool("user_visible", true).
		Msg("Function completed successfully.")

	return result, nil
}

func convertParam(ctx context.Context, mod wasm.Module, schemaType ast.Type, wasmType wasm.ValueType, val any) (uint64, error) {

	switch schemaType.NamedType {

	case "Boolean":
		b, ok := val.(bool)
		if !ok {
			return 0, fmt.Errorf("input value is not a bool")
		}

		// Note, booleans are passed as i32 in wasm
		if wasmType != wasm.ValueTypeI32 {
			return 0, fmt.Errorf("parameter is not defined as a bool on the function")
		}

		if b {
			return 1, nil
		} else {
			return 0, nil
		}

	case "Int":
		n, err := val.(json.Number).Int64()
		if err != nil {
			return 0, fmt.Errorf("input value is not an int")
		}

		switch wasmType {
		case wasm.ValueTypeI32:
			return wasm.EncodeI32(int32(n)), nil
		case wasm.ValueTypeI64:
			return wasm.EncodeI64(n), nil
		default:
			return 0, fmt.Errorf("parameter is not defined as an int on the function")
		}

	case "Float":
		n, err := val.(json.Number).Float64()
		if err != nil {
			return 0, fmt.Errorf("input value is not a float")
		}

		switch wasmType {
		case wasm.ValueTypeF32:
			return wasm.EncodeF32(float32(n)), nil
		case wasm.ValueTypeF64:
			return wasm.EncodeF64(n), nil
		default:
			return 0, fmt.Errorf("parameter is not defined as a float on the function")
		}

	case "String", "ID", "":
		s, ok := val.(string)
		if !ok {
			return 0, fmt.Errorf("input value is not a string")
		}

		// Note, strings are passed as a pointer to a string in wasm memory
		if wasmType != wasm.ValueTypeI32 {
			return 0, fmt.Errorf("parameter is not defined as a string on the function")
		}

		ptr := writeString(ctx, mod, s)
		return uint64(ptr), nil

	default:
		return 0, fmt.Errorf("unknown parameter type: %s", schemaType.NamedType)
	}
}

func convertResult(mem wasm.Memory, schemaType ast.Type, wasmType wasm.ValueType, res uint64) (any, error) {

	switch schemaType.NamedType {

	case "Boolean":
		if wasmType != wasm.ValueTypeI32 {
			return nil, fmt.Errorf("return type is not defined as an bool on the function")
		}

		return res != 0, nil

	case "Int":

		// TODO: Do we need to handle unsigned ints differently?

		switch wasmType {
		case wasm.ValueTypeI32:
			return wasm.DecodeI32(res), nil

		case wasm.ValueTypeI64:
			return int64(res), nil

		default:
			return nil, fmt.Errorf("return type is not defined as an int on the function")
		}

	case "Float":

		switch wasmType {
		case wasm.ValueTypeF32:
			return wasm.DecodeF32(res), nil

		case wasm.ValueTypeF64:
			return wasm.DecodeF64(res), nil

		default:
			return nil, fmt.Errorf("return type is not defined as a float on the function")
		}

	case "ID":
		return nil, fmt.Errorf("the ID scalar is not allowed for function return types (use String instead)")

	default:
		// The return type is either a string, or an object that should be serialized as JSON.
		// Strings are passed as a pointer to a string in wasm memory
		if wasmType != wasm.ValueTypeI32 {
			return nil, fmt.Errorf("return type was not a pointer")
		}

		return readString(mem, uint32(res))
	}
}
