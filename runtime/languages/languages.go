/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package languages

import (
	"fmt"
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

func GetLanguageForSDK(sdk string) (langsupport.Language, error) {

	// strip version if present
	sdkName := sdk
	if i := strings.Index(sdkName, "@"); i != -1 {
		sdkName = sdkName[:i]
	}

	// each SDK has a corresponding language implementation
	switch sdkName {
	case "modus-sdk-as":
		return AssemblyScript(), nil
	case "modus-sdk-go":
		return GoLang(), nil
	}

	return nil, fmt.Errorf("unsupported SDK: %s", sdk)
}
