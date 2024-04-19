/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"fmt"
	"time"

	"hmruntime/functions/assemblyscript"
	"hmruntime/logger"

	"github.com/tetratelabs/wazero/api"
)

func CallFunction(ctx context.Context, mod api.Module, info FunctionInfo, inputs map[string]any) (any, error) {
	fnName := info.Function.Name
	logger.Info(ctx).
		Str("function", fnName).
		Bool("user_visible", true).
		Msg("Calling function.")

	start := time.Now()
	result, err := callFunction(ctx, mod, info, inputs)
	duration := time.Since(start)

	if err != nil {
		err = transformError(err)
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

func callFunction(ctx context.Context, mod api.Module, info FunctionInfo, inputs map[string]any) (any, error) {
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

		param, err := assemblyscript.EncodeValue(ctx, mod, arg.Type, val)
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

	// Handle void result
	if len(res) == 0 {
		return nil, nil
	}

	// Get the result
	result, err := assemblyscript.DecodeValue(ctx, mod, info.Function.ReturnType, res[0])
	if err != nil {
		return nil, fmt.Errorf("function result is invalid: %w", err)
	}

	return result, nil
}
