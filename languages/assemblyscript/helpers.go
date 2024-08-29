/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"hmruntime/utils"
	"reflect"
)

// Allocate memory within the AssemblyScript module.
// This uses the `__new` function exported by the AssemblyScript runtime, so it will be garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func (wa *wasmAdapter) allocateWasmMemory(ctx context.Context, len uint32, classId uint32) (uint32, error) {
	res, err := wa.fnNew.Call(ctx, uint64(len), uint64(classId))
	if err != nil {
		return 0, fmt.Errorf("failed to allocate WASM memory: %w", err)
	}
	return uint32(res[0]), nil
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

func (wa *wasmAdapter) getReflectedType(ctx context.Context, typ string) (reflect.Type, error) {
	if customTypes, ok := ctx.Value(utils.CustomTypesContextKey).(map[string]reflect.Type); ok {
		return wa.typeInfo.getReflectedType(typ, customTypes)
	} else {
		return wa.typeInfo.getReflectedType(typ, nil)
	}
}

func (wa *wasmAdapter) getAssemblyScriptType(ctx context.Context, t reflect.Type) (string, error) {
	if customTypes, ok := ctx.Value(utils.CustomTypesRevContextKey).(map[reflect.Type]string); ok {
		return wa.typeInfo.getAssemblyScriptType(t, customTypes)
	} else {
		return wa.typeInfo.getAssemblyScriptType(t, nil)
	}
}
