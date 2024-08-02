/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"

	wasm "github.com/tetratelabs/wazero/api"
)

// Allocate memory within the AssemblyScript module.
// This uses the `__new` function exported by the AssemblyScript runtime, so it will be garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func (wa *wasmAdapter) allocateWasmMemory(ctx context.Context, mod wasm.Module, len uint32, classId uint32) (uint32, error) {
	fn := mod.ExportedFunction("__new")
	res, err := fn.Call(ctx, uint64(len), uint64(classId))
	if err != nil {
		return 0, fmt.Errorf("failed to allocate WASM memory: %w", err)
	}
	return uint32(res[0]), nil
}

// Pin a managed object in memory within the AssemblyScript module.
// This prevents it from being garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func (wa *wasmAdapter) pinWasmMemory(ctx context.Context, mod wasm.Module, ptr uint32) error {
	fn := mod.ExportedFunction("__pin")
	_, err := fn.Call(ctx, uint64(ptr))
	if err != nil {
		return fmt.Errorf("failed to pin object in WASM memory: %w", err)
	}
	return nil
}

// Unpin a previously-pinned managed object in memory within the AssemblyScript module.
// This allows it to be garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func (wa *wasmAdapter) unpinWasmMemory(ctx context.Context, mod wasm.Module, ptr uint32) error {
	fn := mod.ExportedFunction("__unpin")
	_, err := fn.Call(ctx, uint64(ptr))
	if err != nil {
		return fmt.Errorf("failed to unpin object in WASM memory: %w", err)
	}
	return nil
}

// Sets that arguments length before calling a function that takes a variable number of arguments.
// See https://www.assemblyscript.org/runtime.html#optional-arguments
func (wa *wasmAdapter) setArgumentsLength(ctx context.Context, mod wasm.Module, length int) error {
	fn := mod.ExportedFunction("__setArgumentsLength")
	_, err := fn.Call(ctx, uint64(length))
	if err != nil {
		return fmt.Errorf("failed to set arguments length: %w", err)
	}
	return nil
}
