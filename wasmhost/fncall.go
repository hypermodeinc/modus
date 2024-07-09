/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"fmt"
	"time"

	"hmruntime/functions"
	"hmruntime/functions/assemblyscript"
	"hmruntime/logger"
	"hmruntime/metrics"
	"hmruntime/utils"

	"github.com/rs/xid"
	wasm "github.com/tetratelabs/wazero/api"
)

type ExecutionInfo struct {
	ExecutionId string
	Buffers     *utils.OutputBuffers
	Messages    []utils.LogMessage
	Result      any
}

func CallFunction(ctx context.Context, fnName string, paramValues ...any) (*ExecutionInfo, error) {
	fnInfo, err := functions.GetFunctionInfo(fnName)
	if err != nil {
		return nil, err
	}

	parameters, err := functions.CreateParametersMap(fnInfo, paramValues...)
	if err != nil {
		return nil, err
	}

	return doCallFunction(ctx, fnInfo, parameters)
}

func CallFunctionWithParametersMap(ctx context.Context, fnName string, parameters map[string]any) (*ExecutionInfo, error) {
	fnInfo, err := functions.GetFunctionInfo(fnName)
	if err != nil {
		return nil, err
	}

	return doCallFunction(ctx, fnInfo, parameters)
}

func doCallFunction(ctx context.Context, fnInfo functions.FunctionInfo, parameters map[string]any) (*ExecutionInfo, error) {

	execInfo := ExecutionInfo{
		ExecutionId: xid.New().String(),
		Buffers:     &utils.OutputBuffers{},
		Messages:    []utils.LogMessage{},
	}

	ctx = context.WithValue(ctx, utils.ExecutionIdContextKey, execInfo.ExecutionId)
	ctx = context.WithValue(ctx, utils.FunctionMessagesContextKey, &execInfo.Messages)
	ctx = context.WithValue(ctx, utils.PluginContextKey, fnInfo.Plugin)

	// Each request will get its own instance of the plugin module, so that we can run
	// multiple requests in parallel without risk of corrupting the module's memory.
	// This also protects against security risk, as each request will have its own
	// isolated memory space.  (One request cannot access another request's memory.)

	mod, err := GetModuleInstance(ctx, fnInfo.Plugin, execInfo.Buffers)
	if err != nil {
		return nil, err
	}
	defer mod.Close(ctx)

	logger.Info(ctx).
		Str("function", fnInfo.Function.Name).
		Bool("user_visible", true).
		Msg("Calling function.")

	start := time.Now()
	result, err := invokeFunction(ctx, mod, fnInfo, parameters)
	duration := time.Since(start)

	if err != nil {
		err = functions.TransformError(err)
		logger.Err(ctx, err).
			Str("function", fnInfo.Function.Name).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Error while executing function.")
	} else {
		logger.Info(ctx).
			Str("function", fnInfo.Function.Name).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Function completed successfully.")
	}

	// Update metrics
	metrics.FunctionExecutionsNum.Inc()
	metrics.FunctionExecutionDurationMilliseconds.Observe(float64(duration.Milliseconds()))

	execInfo.Result = result
	return &execInfo, err
}

func invokeFunction(ctx context.Context, mod wasm.Module, info functions.FunctionInfo, parameters map[string]any) (any, error) {

	// Get the wasm function
	fn := mod.ExportedFunction(info.Function.Name)
	if fn == nil {
		return nil, fmt.Errorf("function %s not found in plugin %s", info.Function.Name, info.Plugin.Name())
	}

	// Get parameters to pass as input to the plugin function
	params := make([]uint64, len(info.Function.Parameters))
	param_mask := uint64(0)

	has_opt_param := false
	// Set all bits to 0
	for i, arg := range info.Function.Parameters {
		val := parameters[arg.Name]

		if arg.Optional {
			has_opt_param = true
			if val == nil {
				param_mask &= ^(1 << i)
				continue
			} else {
				param_mask |= 1 << i
			}
		}

		if val == nil {
			return nil, fmt.Errorf("parameter '%s' is missing", arg.Name)
		}

		param, err := assemblyscript.EncodeValue(ctx, mod, arg.Type, val)
		if err != nil {
			return nil, fmt.Errorf("function parameter '%s' is invalid: %w", arg.Name, err)
		}

		params[i] = param
	}
	if has_opt_param {
		params = append(params, param_mask)
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
