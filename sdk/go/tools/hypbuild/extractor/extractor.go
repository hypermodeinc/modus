/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package extractor

import (
	"go/types"
	"sort"

	"github.com/hypermodeAI/functions-go/tools/hypbuild/config"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/metadata"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/utils"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/wasm"
)

func CollectProgramInfo(config *config.Config, meta *metadata.Metadata, wasmFunctions *wasm.WasmFunctions) error {
	pkgs, err := loadPackages(config.SourceDir)
	if err != nil {
		return err
	}

	requiredTypes := make(map[string]types.Type)

	for name, f := range getExportedFunctions(pkgs) {
		if _, ok := wasmFunctions.Exports[name]; ok {
			meta.FnExports[name] = transformFunc(name, f)
			findRequiredTypes(f, requiredTypes)
		}
	}

	for name, f := range getImportedFunctions(pkgs) {
		if _, ok := wasmFunctions.Imports[name]; ok {
			meta.FnImports[name] = transformFunc(name, f)
			findRequiredTypes(f, requiredTypes)
		}
	}

	// proxy imports overwrite regular imports
	for name, f := range getProxyImportFunctions(pkgs) {
		if _, ok := meta.FnImports[name]; ok {
			meta.FnImports[name] = transformFunc(name, f)
			findRequiredTypes(f, requiredTypes)
		}
	}

	id := uint32(3) // 1 and 2 are reserved for []byte and string
	keys := utils.MapKeys(requiredTypes)
	sort.Strings(keys)
	for _, name := range keys {
		t := requiredTypes[name]
		if s, ok := t.(*types.Struct); ok && !wellKnownTypes[name] {
			t := transformStruct(name, s)
			t.Id = id
			meta.Types[name] = t
		} else {
			meta.Types[name] = &metadata.TypeDefinition{
				Id:   id,
				Name: name,
			}
		}
		id++
	}

	return nil
}
