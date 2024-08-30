/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"
	"reflect"

	"hypruntime/utils"
)

func (wa *wasmAdapter) allocateWasmMemory(ctx context.Context, size uint32) (uint32, utils.Cleaner, error) {
	res, err := wa.fnMalloc.Call(ctx, uint64(size))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to allocate WASM memory: %w", err)
	}
	ptr := uint32(res[0])

	if ptr == 0 {
		return 0, nil, fmt.Errorf("failed to allocate WASM memory")
	}

	cln := utils.NewCleanerN(1)
	cln.AddCleanup(func() error {
		if _, err := wa.fnFree.Call(ctx, uint64(ptr)); err != nil {
			return fmt.Errorf("failed to free WASM memory: %w", err)
		}
		return nil
	})

	return ptr, cln, nil
}

func (wa *wasmAdapter) newWasmObject(ctx context.Context, id uint32) (uint32, utils.Cleaner, error) {
	res, err := wa.fnNew.Call(ctx, uint64(id))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create an object with id %d in WASM memory: %w", id, err)
	}
	ptr := uint32(res[0])

	if ptr == 0 {
		return 0, nil, fmt.Errorf("failed to create an object with id %d in WASM memory", id)
	}

	cln := utils.NewCleanerN(1)
	cln.AddCleanup(func() error {
		if _, err := wa.fnUnpin.Call(ctx, uint64(ptr)); err != nil {
			return fmt.Errorf("failed to unpin WASM object: %w", err)
		}
		return nil
	})

	return ptr, cln, nil
}

func (wa *wasmAdapter) makeWasmObject(ctx context.Context, id, size uint32) (uint32, utils.Cleaner, error) {
	res, err := wa.fnMake.Call(ctx, uint64(id), uint64(size))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to make an object with id %d and size %d in WASM memory: %w", id, size, err)
	}
	ptr := uint32(res[0])

	if ptr == 0 {
		return 0, nil, fmt.Errorf("failed to make an object with id %d and size %d in WASM memory", id, size)
	}

	cln := utils.NewCleanerN(1)
	cln.AddCleanup(func() error {
		if _, err := wa.fnUnpin.Call(ctx, uint64(ptr)); err != nil {
			return fmt.Errorf("failed to unpin WASM object: %w", err)
		}
		return nil
	})

	return ptr, cln, nil
}

func (wa *wasmAdapter) getReflectedType(ctx context.Context, typ string) (reflect.Type, error) {
	if customTypes, ok := ctx.Value(utils.CustomTypesContextKey).(map[string]reflect.Type); ok {
		return wa.typeInfo.getReflectedType(typ, customTypes)
	} else {
		return wa.typeInfo.getReflectedType(typ, nil)
	}
}
