/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"errors"

	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

type param struct {
	ptr uint32
	val any
}

func readParams(ctx context.Context, mod wasm.Module, params ...param) error {
	plugin := ctx.Value(utils.PluginContextKey).(*plugins.Plugin)
	adapter := plugin.Language.WasmAdapter()

	errs := make([]error, 0, len(params))
	for _, p := range params {
		if err := adapter.DecodeValue(ctx, mod, uint64(p.ptr), p.val); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func writeResult(ctx context.Context, mod wasm.Module, val any) (uint32, error) {
	plugin := ctx.Value(utils.PluginContextKey).(*plugins.Plugin)
	adapter := plugin.Language.WasmAdapter()

	p, err := adapter.EncodeValue(ctx, mod, val)
	return uint32(p), err
}
