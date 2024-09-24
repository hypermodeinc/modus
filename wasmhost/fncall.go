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

	"hypruntime/functions"
	"hypruntime/logger"
	"hypruntime/metrics"
	"hypruntime/utils"

	"github.com/rs/xid"
	"github.com/tetratelabs/wazero/sys"
)

type ExecutionInfo interface {
	ExecutionId() string
	Buffers() utils.OutputBuffers
	Messages() []utils.LogMessage
	Result() any
}

type executionInfo struct {
	executionId string
	buffers     utils.OutputBuffers
	messages    []utils.LogMessage
	result      any
}

func (e *executionInfo) ExecutionId() string {
	return e.executionId
}

func (e *executionInfo) Buffers() utils.OutputBuffers {
	return e.buffers
}

func (e *executionInfo) Messages() []utils.LogMessage {
	return e.messages
}

func (e *executionInfo) Result() any {
	return e.result
}

func CallFunction(ctx context.Context, fnName string, paramValues ...any) (ExecutionInfo, error) {
	return GetWasmHost(ctx).CallFunctionByName(ctx, fnName, paramValues...)
}

func (host *wasmHost) CallFunctionByName(ctx context.Context, fnName string, paramValues ...any) (ExecutionInfo, error) {
	info, err := host.GetFunctionInfo(fnName)
	if err != nil {
		return nil, err
	}

	fn := info.Metadata()
	parameters, err := functions.CreateParametersMap(fn, paramValues...)
	if err != nil {
		return nil, err
	}

	return host.CallFunction(ctx, info, parameters)
}

func (host *wasmHost) CallFunction(ctx context.Context, fnInfo functions.FunctionInfo, parameters map[string]any) (ExecutionInfo, error) {

	execInfo := &executionInfo{
		executionId: xid.New().String(),
		buffers:     utils.NewOutputBuffers(),
		messages:    []utils.LogMessage{},
	}

	fnName := fnInfo.Name()
	plugin := fnInfo.Plugin()
	plan := fnInfo.ExecutionPlan()

	ctx = context.WithValue(ctx, utils.ExecutionIdContextKey, execInfo.executionId)
	ctx = context.WithValue(ctx, utils.FunctionMessagesContextKey, &execInfo.messages)
	ctx = context.WithValue(ctx, utils.FunctionNameContextKey, fnName)
	ctx = context.WithValue(ctx, utils.PluginContextKey, plugin)
	ctx = context.WithValue(ctx, utils.MetadataContextKey, plugin.Metadata)
	ctx = context.WithValue(ctx, utils.WasmHostContextKey, host)

	// Each request will get its own instance of the plugin module, so that we can run
	// multiple requests in parallel without risk of corrupting the module's memory.
	// This also protects against security risk, as each request will have its own
	// isolated memory space.  (One request cannot access another request's memory.)

	mod, err := host.GetModuleInstance(ctx, plugin, execInfo.buffers)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting module instance.")
		return nil, err
	}
	defer mod.Close(ctx)

	wa := plugin.Language.NewWasmAdapter(mod)
	ctx = context.WithValue(ctx, utils.WasmAdapterContextKey, wa)

	logger.Info(ctx).
		Str("function", fnName).
		Bool("user_visible", true).
		Msg("Calling function.")

	start := time.Now()
	result, err := plan.InvokeFunction(ctx, wa, parameters)
	duration := time.Since(start)

	exitErr := &sys.ExitError{}

	if err == nil {
		logger.Info(ctx).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Function completed successfully.")
	} else if errors.As(err, &exitErr) {
		exitCode := int32(exitErr.ExitCode())
		logger.Error(ctx).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Int32("exit_code", exitCode).
			Msgf("Function ended prematurely with exit code %d.", exitCode)
	} else if errors.Is(err, context.Canceled) {
		logger.Warn(ctx).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Function execution was canceled.")
	} else {
		if utils.HypermodeDebugEnabled() {
			fmt.Fprintln(os.Stderr, err)
		}
		logger.Err(ctx, err).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Error while executing function.")
	}

	// Update metrics
	metrics.FunctionExecutionsNum.Inc()
	d := float64(duration.Milliseconds())
	metrics.FunctionExecutionDurationMilliseconds.Observe(d)
	metrics.FunctionExecutionDurationMillisecondsSummary.Observe(d)

	execInfo.result = result
	return execInfo, err
}
