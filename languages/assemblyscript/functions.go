/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"

	"hmruntime/plugins/metadata"
)

func (wa *wasmAdapter) InvokeFunction(ctx context.Context, function *metadata.Function, parameters map[string]any) (result any, err error) {

	// Get the wasm function
	fn := wa.mod.ExportedFunction(function.Name)
	if fn == nil {
		return nil, fmt.Errorf("function %s not found in wasm module", function.Name)
	}

	// Get parameters to pass as input to the function
	params, has_defaults, err := wa.GetParameters(ctx, function.Parameters, parameters)
	if err != nil {
		return nil, err
	}

	// If the function has any parameters with default values, we need to set the arguments length.
	// Since we pass all the arguments ourselves, we just need to set the total length of the arguments.
	if has_defaults {
		err := wa.setArgumentsLength(ctx, len(params))
		if err != nil {
			return nil, err
		}
	}

	// Call the function
	res, err := fn.Call(ctx, params...)
	if err != nil {
		return nil, err
	}

	// Handle void result
	if len(res) == 0 {
		return nil, nil
	}

	// Get the result
	resultType := function.Results[0].Type
	if result, err := wa.decodeValue(ctx, resultType, res[0]); err != nil {
		return nil, fmt.Errorf("function result is invalid: %w", err)
	} else {
		return result, nil
	}
}

func (wa *wasmAdapter) GetParameters(ctx context.Context, paramInfo []*metadata.Parameter, parameters map[string]any) ([]uint64, bool, error) {
	params := make([]uint64, len(paramInfo))
	mask := uint64(0)
	has_opt := false
	has_def := false

	for i, arg := range paramInfo {
		val, found := parameters[arg.Name]

		if arg.Default != nil {
			has_def = true
			if !found {
				val = *arg.Default
			}
		}

		// maintain compatibility with the deprecated "optional" field
		if arg.Optional {
			has_opt = true
			if !found {
				continue
			}
		}

		if val == nil {
			if wa.typeInfo.IsNullable(arg.Type) {
				continue
			}
			return nil, false, fmt.Errorf("parameter '%s' cannot be null", arg.Name)
		}

		mask |= 1 << i

		param, err := wa.encodeValueForParameter(ctx, arg.Type, val)
		if err != nil {
			return nil, false, fmt.Errorf("function parameter '%s' is invalid: %w", arg.Name, err)
		}

		params[i] = param
	}

	if has_opt {
		params = append(params, mask)
	}

	return params, has_def, nil
}
