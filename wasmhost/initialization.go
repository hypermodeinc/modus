/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"

	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
)

func InitWasmHost(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	cfg := wazero.NewRuntimeConfig()
	cfg = cfg.WithCloseOnContextDone(true)
	RuntimeInstance = wazero.NewRuntimeWithConfig(ctx, cfg)
}
