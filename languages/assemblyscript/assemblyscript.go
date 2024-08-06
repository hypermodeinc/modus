/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

const LanguageName = "AssemblyScript"

var _typeInfoProvider = &typeInfoProvider{}

var _wasmAdapter = &wasmAdapter{
	typeInfo: _typeInfoProvider,
}

func TypeInfo() *typeInfoProvider {
	return _typeInfoProvider
}

func WasmAdapter() *wasmAdapter {
	return _wasmAdapter
}

type typeInfoProvider struct{}

type wasmAdapter struct {
	typeInfo *typeInfoProvider
}
