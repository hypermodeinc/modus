/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	wasm "github.com/tetratelabs/wazero/api"
)

const LanguageName = "Go"

var _typeInfoProvider = &typeInfoProvider{}

func TypeInfo() *typeInfoProvider {
	return _typeInfoProvider
}

func NewWasmAdapter(mod wasm.Module) *wasmAdapter {
	return &wasmAdapter{
		mod:         mod,
		typeInfo:    _typeInfoProvider,
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

type typeInfoProvider struct{}

type wasmAdapter struct {
	mod         wasm.Module
	typeInfo    *typeInfoProvider
	visitedPtrs map[uint32]int
	fnMalloc    wasm.Function
	fnFree      wasm.Function
	fnNew       wasm.Function
	fnMake      wasm.Function
	fnUnpin     wasm.Function
	fnReadMap   wasm.Function
	fnWriteMap  wasm.Function
}
