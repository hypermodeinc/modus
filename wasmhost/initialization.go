/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"

	"hmruntime/logger"
	"hmruntime/utils"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
)

func InitWasmHost(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	configureLogger()

	cfg := wazero.NewRuntimeConfig()
	cfg = cfg.WithCloseOnContextDone(true)
	RuntimeInstance = wazero.NewRuntimeWithConfig(ctx, cfg)
}

func configureLogger() {
	logger.AddAdapter(func(ctx context.Context, lc zerolog.Context) zerolog.Context {
		if executionId, ok := ctx.Value(utils.ExecutionIdContextKey).(string); ok {
			lc = lc.Str("execution_id", executionId)
		}

		return lc
	})
}
