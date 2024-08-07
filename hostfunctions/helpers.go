/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hmruntime/languages"
	"hmruntime/logger"
	"hmruntime/plugins"

	wasm "github.com/tetratelabs/wazero/api"
)

func readParams(ctx context.Context, mod wasm.Module, stack []uint64, params ...any) error {
	if len(params) != len(stack) {
		return fmt.Errorf("expected %d arguments, got %d", len(params), len(stack))
	}

	adapter, err := getWasmAdapter(ctx)
	if err != nil {
		return err
	}

	errs := make([]error, 0, len(params))
	for i, p := range params {
		if err := adapter.DecodeValue(ctx, mod, stack[i], p); err != nil {
			errs = append(errs, err)
		}
		stack[i] = 0 // clear the value from the stack
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func writeResults(ctx context.Context, mod wasm.Module, stack []uint64, results ...any) error {
	if len(results) != len(stack) {
		return fmt.Errorf("expected %d results, got %d", len(results), len(stack))
	}

	adapter, err := getWasmAdapter(ctx)
	if err != nil {
		return err
	}

	errs := make([]error, 0, len(results))
	for i, r := range results {
		val, err := adapter.EncodeValue(ctx, mod, r)
		if err != nil {
			errs = append(errs, err)
		}
		stack[i] = val
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil

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

type hostFunctionDefinition struct {
	name     string
	function wasm.GoModuleFunction
	params   []wasm.ValueType
	results  []wasm.ValueType
}

// Each message is optional, but if provided, it will be logged at the appropriate time.
type hostFunctionMessages struct {
	Starting  string
	Completed string
	Cancelled string
	Error     string
	Detail    string
}

func callHostFunction(ctx context.Context, fn func() error, msgs *hostFunctionMessages) bool {

	if msgs != nil && msgs.Starting != "" {
		l := logger.Info(ctx).Bool("user_visible", true)
		if msgs.Detail != "" {
			l.Str("detail", msgs.Detail)
		}
		l.Msg(msgs.Starting)
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	if errors.Is(err, context.Canceled) {
		if msgs != nil && msgs.Cancelled != "" {
			l := logger.Warn(ctx).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.Detail != "" {
				l.Str("detail", msgs.Detail)
			}
			l.Msg(msgs.Cancelled)
		}
		return false
	} else if err != nil {
		if msgs != nil && msgs.Error != "" {
			l := logger.Err(ctx, err).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.Detail != "" {
				l.Str("detail", msgs.Detail)
			}
			l.Msg(msgs.Error)
		}
		return false
	} else {
		if msgs != nil && msgs.Completed != "" {
			l := logger.Info(ctx).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.Detail != "" {
				l.Str("detail", msgs.Detail)
			}
			l.Msg(msgs.Completed)
		}
		return true
	}
}
