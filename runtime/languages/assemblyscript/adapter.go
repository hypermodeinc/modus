/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript

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
		mod:                  mod,
		visitedPtrs:          make(map[uint32]int),
		fnNew:                mod.ExportedFunction("__new"),
		fnPin:                mod.ExportedFunction("__pin"),
		fnUnpin:              mod.ExportedFunction("__unpin"),
		fnSetArgumentsLength: mod.ExportedFunction("__setArgumentsLength"),
	}
}

// for testing purposes only
func NewWasmAdapter_internal(mod wasm.Module) *wasmAdapter {
	return NewWasmAdapter(mod).(*wasmAdapter)
}

type wasmAdapter struct {
	mod                  wasm.Module
	visitedPtrs          map[uint32]int
	fnNew                wasm.Function
	fnPin                wasm.Function
	fnUnpin              wasm.Function
	fnSetArgumentsLength wasm.Function
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
	// If the function has any parameters with default values, we need to set the arguments length before invoking the functions.
	// Since we pass all the arguments ourselves, we just need to set the total length of the arguments.
	if plan.HasDefaultParameters() {
		return wa.setArgumentsLength(ctx, len(plan.FnMetadata().Parameters))
	}
	return nil
}

func (wa *wasmAdapter) AllocateMemory(ctx context.Context, size uint32) (uint32, utils.Cleaner, error) {
	// Array buffers (id:1) are used for general memory allocation in AssemblyScript.
	return wa.allocateAndPinMemory(ctx, size, 1)
}

// Allocate and pin memory within the AssemblyScript module.
// The cleaner returned will unpin the memory when invoked.
func (wa *wasmAdapter) allocateAndPinMemory(ctx context.Context, size, classId uint32) (uint32, utils.Cleaner, error) {
	ptr, err := wa.allocateWasmMemory(ctx, size, classId)
	if err != nil {
		return 0, nil, err
	}

	if err := wa.pinWasmMemory(ctx, ptr); err != nil {
		return 0, nil, err
	}

	cln := utils.NewCleanerN(1)
	cln.AddCleanup(func() error {
		return wa.unpinWasmMemory(ctx, ptr)
	})

	return ptr, cln, nil
}

// Allocate memory within the AssemblyScript module.
// This uses the `__new` function exported by the AssemblyScript runtime, so it will be garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func (wa *wasmAdapter) allocateWasmMemory(ctx context.Context, size, classId uint32) (uint32, error) {
	res, err := wa.fnNew.Call(ctx, uint64(size), uint64(classId))
	if err != nil {
		return 0, fmt.Errorf("failed to allocate WASM memory (size: %d, id: %d): %w", size, classId, err)
	}

	ptr := uint32(res[0])
	if ptr == 0 {
		return 0, errors.New("failed to allocate WASM memory")
	}

	return ptr, nil
}

// Pin a managed object in memory within the AssemblyScript module.
// This prevents it from being garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func (wa *wasmAdapter) pinWasmMemory(ctx context.Context, ptr uint32) error {
	_, err := wa.fnPin.Call(ctx, uint64(ptr))
	if err != nil {
		return fmt.Errorf("failed to pin object in WASM memory: %w", err)
	}
	return nil
}

// Unpin a previously-pinned managed object in memory within the AssemblyScript module.
// This allows it to be garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func (wa *wasmAdapter) unpinWasmMemory(ctx context.Context, ptr uint32) error {
	_, err := wa.fnUnpin.Call(ctx, uint64(ptr))
	if err != nil {
		return fmt.Errorf("failed to unpin object in WASM memory: %w", err)
	}
	return nil
}

// Sets that arguments length before calling a function that takes a variable number of arguments.
// See https://www.assemblyscript.org/runtime.html#optional-arguments
func (wa *wasmAdapter) setArgumentsLength(ctx context.Context, length int) error {
	_, err := wa.fnSetArgumentsLength.Call(ctx, uint64(length))
	if err != nil {
		return fmt.Errorf("failed to set arguments length: %w", err)
	}
	return nil
}
