/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package wasmhost

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hypermodeinc/modus/runtime/functions"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/metrics"
	"github.com/hypermodeinc/modus/runtime/utils"

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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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
		// NOTE: This can occur if the function calls `exit` or `abort` in the WASM code, or if they throw an exception or panic.
		// In those cases, the message of the exception or panic will have already been logged via the `log` host function.
		exitCode := int32(exitErr.ExitCode())
		if exitCode == 0 {
			// exit code 0 is a non-error exit, so we log it as a warning
			logger.Warn(ctx).
				Str("function", fnName).
				Dur("duration_ms", duration).
				Bool("user_visible", true).
				Int32("exit_code", exitCode).
				Msgf("Function ended prematurely with exit code %d.  It is generally better to return properly than to exit abruptly.", exitCode)
		} else {
			// any other exit code is an error exit, so we log it as an error
			logger.Error(ctx).
				Str("function", fnName).
				Dur("duration_ms", duration).
				Bool("user_visible", true).
				Int32("exit_code", exitCode).
				Msgf("Function ended prematurely with exit code %d.  This may have been intentional, or caused by an exception or panic in your code.", exitCode)
		}
	} else if errors.Is(err, context.Canceled) {
		// Cancellation is not an error, but we still want to log it.
		// This can occur if the function takes too long to execute, or if the user cancels the request.
		logger.Warn(ctx).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("Function execution was canceled.")
	} else {
		// While debugging, it helps if we can see the error in the console without escaped newlines and other json formatting.
		if utils.DebugModeEnabled() {
			fmt.Fprintln(os.Stderr, err)
		}
		// NOTE: Errors of this type should not be user-visible, as they were caused by some Runtime issue, not the user's code.
		// This will also ensure the error is reported to Sentry.
		logger.Err(ctx, err).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Msg("Error while executing function.")

		// However, we should still log _something_ that is user visible, so that the user knows something went wrong when they look at the function run logs.
		logger.Error(ctx).
			Str("function", fnName).
			Dur("duration_ms", duration).
			Bool("user_visible", true).
			Msg("An internal runtime error occurred while executing the function.")
	}

	// Update metrics
	metrics.FunctionExecutionsNum.Inc()
	d := float64(duration.Milliseconds())
	metrics.FunctionExecutionDurationMilliseconds.WithLabelValues(fnName).Observe(d)
	metrics.FunctionExecutionDurationMillisecondsSummary.WithLabelValues(fnName).Observe(d)

	execInfo.result = result
	return execInfo, err
}
