/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"errors"

	"hmruntime/languages"
	"hmruntime/plugins"

	wasm "github.com/tetratelabs/wazero/api"
)

type param struct {
	ptr uint32
	val any
}

func readParams(ctx context.Context, mod wasm.Module, params ...param) error {
	adapter, err := getWasmAdapter(ctx)
	if err != nil {
		return err
	}

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
	adapter, err := getWasmAdapter(ctx)
	if err != nil {
		return 0, err
	}

	p, err := adapter.EncodeValue(ctx, mod, val)
	return uint32(p), err
}

func getWasmAdapter(ctx context.Context) (languages.WasmAdapter, error) {
	p := plugins.GetPlugin(ctx)
	if p == nil {
		return nil, errors.New("no plugin found in context")
	}

	wa := p.Language.WasmAdapter()
	if wa == nil {
		return nil, errors.New("no wasm adapter found in plugin")
	}

	return wa, nil
}
