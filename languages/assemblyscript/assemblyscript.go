/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	wasm "github.com/tetratelabs/wazero/api"
)

const LanguageName = "AssemblyScript"

var _typeInfoProvider = &typeInfoProvider{}

func TypeInfo() *typeInfoProvider {
	return _typeInfoProvider
}

func NewWasmAdapter(mod wasm.Module) *wasmAdapter {
	return &wasmAdapter{
		mod:                  mod,
		typeInfo:             _typeInfoProvider,
		visitedPtrs:          make(map[uint32]int),
		fnNew:                mod.ExportedFunction("__new"),
		fnPin:                mod.ExportedFunction("__pin"),
		fnUnpin:              mod.ExportedFunction("__unpin"),
		fnSetArgumentsLength: mod.ExportedFunction("__setArgumentsLength"),
	}
}

type typeInfoProvider struct{}

type wasmAdapter struct {
	mod                  wasm.Module
	typeInfo             *typeInfoProvider
	visitedPtrs          map[uint32]int
	fnNew                wasm.Function
	fnPin                wasm.Function
	fnUnpin              wasm.Function
	fnSetArgumentsLength wasm.Function
}
