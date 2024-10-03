/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang

import (
	"context"
	"errors"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func NewWasmAdapter(mod wasm.Module) langsupport.WasmAdapter {
	return &wasmAdapter{
		mod:         mod,
		visitedPtrs: make(map[uint32]int),
		fnMalloc:    mod.ExportedFunction("malloc"),
		fnFree:      mod.ExportedFunction("free"),
		fnNew:       mod.ExportedFunction("__new"),
		fnMake:      mod.ExportedFunction("__make"),
		fnUnpin:     mod.ExportedFunction("__unpin"),
		fnReadMap:   mod.ExportedFunction("__read_map"),
		fnWriteMap:  mod.ExportedFunction("__write_map"),
	}
}

type wasmAdapter struct {
	mod         wasm.Module
	visitedPtrs map[uint32]int
	fnMalloc    wasm.Function
	fnFree      wasm.Function
	fnNew       wasm.Function
	fnMake      wasm.Function
	fnUnpin     wasm.Function
	fnReadMap   wasm.Function
	fnWriteMap  wasm.Function
}

func (*wasmAdapter) TypeInfo() langsupport.LanguageTypeInfo {
	return _langTypeInfo
}

func (wa *wasmAdapter) Memory() wasm.Memory {
	return wa.mod.Memory()
}

func (wa *wasmAdapter) GetFunction(name string) wasm.Function {
	return wa.mod.ExportedFunction(name)
}

func (wa *wasmAdapter) PreInvoke(ctx context.Context, plan langsupport.ExecutionPlan) error {
	return nil
}

func (wa *wasmAdapter) AllocateMemory(ctx context.Context, size uint32) (uint32, utils.Cleaner, error) {
	res, err := wa.fnMalloc.Call(ctx, uint64(size))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to allocate WASM memory (size: %d): %w", size, err)
	}

	ptr := uint32(res[0])
	if ptr == 0 {
		return 0, nil, errors.New("failed to allocate WASM memory")
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
