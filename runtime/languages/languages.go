/*
 * Copyright 2024 Hypermode, Inc.
 */

package languages

import (
	"strings"

	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/languages/assemblyscript"
	"github.com/hypermodeinc/modus/runtime/languages/golang"
)

var lang_AssemblyScript = langsupport.NewLanguage(
	"AssemblyScript",
	assemblyscript.LanguageTypeInfo(),
	assemblyscript.NewPlanner,
	assemblyscript.NewWasmAdapter,
)

var lang_Go = langsupport.NewLanguage(
	"Go",
	golang.LanguageTypeInfo(),
	golang.NewPlanner,
	golang.NewWasmAdapter,
)

func AssemblyScript() langsupport.Language {
	return lang_AssemblyScript
}

func GoLang() langsupport.Language {
	return lang_Go
}

func GetLanguageForSDK(sdk string) langsupport.Language {

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
