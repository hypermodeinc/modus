/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"

	"hmruntime/logger"
	"hmruntime/utils"

	"github.com/rs/zerolog"
)

// TODO: refactor to remove global

var GlobalWasmHost *WasmHost

func InitWasmHost(ctx context.Context, opts ...func(*WasmHost) error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	configureLogger()

	GlobalWasmHost = NewWasmHost(ctx, opts...)
}

func configureLogger() {
	logger.AddAdapter(func(ctx context.Context, lc zerolog.Context) zerolog.Context {
		if executionId, ok := ctx.Value(utils.ExecutionIdContextKey).(string); ok {
			lc = lc.Str("execution_id", executionId)
		}

		return lc
	})
}
