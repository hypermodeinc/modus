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
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/rs/zerolog"
)

func InitWasmHost(ctx context.Context, registrations ...func(WasmHost) error) WasmHost {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	configureLogger()

	return NewWasmHost(ctx, registrations...)
}

func configureLogger() {
	logger.AddAdapter(func(ctx context.Context, lc zerolog.Context) zerolog.Context {
		if executionId, ok := ctx.Value(utils.ExecutionIdContextKey).(string); ok {
			lc = lc.Str("execution_id", executionId)
		}

		return lc
	})
}
