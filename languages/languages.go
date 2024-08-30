/*
 * Copyright 2024 Hypermode, Inc.
 */

package languages

import (
	"strings"

	"hypruntime/languages/assemblyscript"
	"hypruntime/languages/golang"

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

var lang_Go = &language{
	name:     golang.LanguageName,
	typeInfo: golang.TypeInfo(),
	wasmAdapterFactory: func(mod wasm.Module) WasmAdapter {
		return golang.NewWasmAdapter(mod)
	},
}

func GoLang() Language {
	return lang_Go
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
	case "functions-go":
		return GoLang()
	}

	return nil
}
