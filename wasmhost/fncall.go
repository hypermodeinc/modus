/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"hmruntime/functions"
	"hmruntime/logger"
	"hmruntime/metrics"
	"hmruntime/plugins"
	"hmruntime/plugins/metadata"
	"hmruntime/utils"

	"github.com/rs/xid"
	"github.com/tetratelabs/wazero/sys"
)

type ExecutionInfo struct {
	ExecutionId string
	Buffers     *utils.OutputBuffers
	Messages    []utils.LogMessage
	Result      any
}

func CallFunction(ctx context.Context, fnName string, paramValues ...any) (*ExecutionInfo, error) {
	return GlobalWasmHost.CallFunction(ctx, fnName, paramValues...)
}

func (host *WasmHost) CallFunction(ctx context.Context, fnName string, paramValues ...any) (*ExecutionInfo, error) {
	function, plugin, err := functions.GetFunctionAndPlugin(fnName)
	if err != nil {
		return nil, err
	}

	parameters, err := functions.CreateParametersMap(function, paramValues...)
	if err != nil {
		return nil, err
	}

	return host.InvokeFunction(ctx, plugin, function, parameters)
}

func (host *WasmHost) CallFunctionWithParametersMap(ctx context.Context, fnName string, parameters map[string]any) (*ExecutionInfo, error) {
	function, plugin, err := functions.GetFunctionAndPlugin(fnName)
	if err != nil {
		return nil, err
	}

	return host.InvokeFunction(ctx, plugin, function, parameters)
}

func (host *WasmHost) InvokeFunction(ctx context.Context, plugin *plugins.Plugin, function *metadata.Function, parameters map[string]any) (*ExecutionInfo, error) {

	execInfo := ExecutionInfo{
		ExecutionId: xid.New().String(),
		Buffers:     &utils.OutputBuffers{},
		Messages:    []utils.LogMessage{},
	}

	ctx = context.WithValue(ctx, utils.ExecutionIdContextKey, execInfo.ExecutionId)
	ctx = context.WithValue(ctx, utils.FunctionMessagesContextKey, &execInfo.Messages)
	ctx = context.WithValue(ctx, utils.FunctionNameContextKey, function.Name)
	ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)
	ctx = context.WithValue(ctx, utils.MetadataContextKey, plugin.Metadata)

	// Each request will get its own instance of the plugin module, so that we can run
	// multiple requests in parallel without risk of corrupting the module's memory.
	// This also protects against security risk, as each request will have its own
	// isolated memory space.  (One request cannot access another request's memory.)

	mod, err := host.GetModuleInstance(ctx, plugin, execInfo.Buffers)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting module instance.")
		return nil, err
	}
	defer mod.Close(ctx)

	wa := plugin.Language.NewWasmAdapter(mod)
	ctx = context.WithValue(ctx, utils.WasmAdapterContextKey, wa)

	logger.Info(ctx).
		Str("function", function.Name).
		Bool("user_visible", true).
		Msg("Calling function.")

	start := time.Now()
	result, err := func() (result any, err error) {
		defer func() {
			if e := utils.ConvertToError(recover()); e != nil {
				err = e
			}
		}()
		return wa.InvokeFunction(ctx, function, parameters)
	}()
	duration := time.Since(start)

	exitErr := &sys.ExitError{}

	if err == nil {
		logger.Info(ctx).
			Str("function", function.Name).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Function completed successfully.")
	} else if errors.As(err, &exitErr) {
		exitCode := int32(exitErr.ExitCode())
		logger.Error(ctx).
			Str("function", function.Name).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Int32("exit_code", exitCode).
			Msgf("Function ended prematurely with exit code %d.", exitCode)
	} else if errors.Is(err, context.Canceled) {
		logger.Warn(ctx).
			Str("function", function.Name).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Function execution was canceled.")
	} else {
		if utils.HypermodeDebugEnabled() {
			fmt.Fprintln(os.Stderr, err)
		}
		logger.Err(ctx, err).
			Str("function", function.Name).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Error while executing function.")
	}

	// Update metrics
	metrics.FunctionExecutionsNum.Inc()
	metrics.FunctionExecutionDurationMilliseconds.Observe(float64(duration.Milliseconds()))

	execInfo.Result = result
	return &execInfo, err
}
