/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

import (
	"context"
	"errors"

	"hypruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

type WasmAdapter interface {
	TypeInfo() LanguageTypeInfo
	Memory() wasm.Memory
	AllocateMemory(ctx context.Context, size uint32) (uint32, utils.Cleaner, error)
	GetFunction(name string) wasm.Function
	PreInvoke(ctx context.Context, plan ExecutionPlan) error
}

func GetWasmAdapter(ctx context.Context) (WasmAdapter, error) {
	if wa, ok := ctx.Value(utils.WasmAdapterContextKey).(WasmAdapter); ok {
		return wa, nil
	}
	return nil, errors.New("no WasmAdapter in context")
}
