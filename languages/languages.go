/*
 * Copyright 2024 Hypermode, Inc.
 */

package languages

import (
	"strings"

	"hmruntime/languages/assemblyscript"

	wasm "github.com/tetratelabs/wazero/api"
)

var lang_AssemblyScript = &language{
	name:     assemblyscript.LanguageName,
	typeInfo: assemblyscript.TypeInfo(),
	wasmAdapterFactory: func(mod wasm.Module) WasmAdapter {
		return assemblyscript.NewWasmAdapter(mod)
	},
}

func AssemblyScript() Language {
	return lang_AssemblyScript
}

func GetLanguageForSDK(sdk string) Language {

	// strip version if present
	if i := strings.Index(sdk, "@"); i != -1 {
		sdk = sdk[:i]
	}

	// each SDK has a corresponding language implementation
	switch sdk {
	case "functions-as":
		return AssemblyScript()
	}

	return nil
}
