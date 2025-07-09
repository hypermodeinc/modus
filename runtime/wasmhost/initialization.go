/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package wasmhost

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/rs/zerolog"
)

func InitWasmHost(ctx context.Context, registrations ...func(WasmHost) error) WasmHost {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	configureLogger()

	return NewWasmHost(ctx, registrations...)
}

func configureLogger() {
	logger.AddAdapter(func(ctx context.Context, lc zerolog.Context) zerolog.Context {
		if executionId, ok := ctx.Value(utils.ExecutionIdContextKey).(string); ok {
			lc = lc.Str("execution_id", executionId)
		}

		if fnName, ok := ctx.Value(utils.FunctionNameContextKey).(string); ok {
			lc = lc.Str("function", fnName)
		}

		if agentName, ok := ctx.Value(utils.AgentNameContextKey).(string); ok {
			lc = lc.Str("agent", agentName)
		}

		if agentId, ok := ctx.Value(utils.AgentIdContextKey).(string); ok {
			lc = lc.Str("agent_id", agentId)
		}

		return lc
	})
}
