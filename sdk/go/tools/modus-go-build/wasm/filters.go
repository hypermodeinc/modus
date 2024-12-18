/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package wasm

import (
	"path/filepath"

	"github.com/hypermodeinc/modus/lib/wasmextractor"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/config"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metadata"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
)

func FilterMetadata(config *config.Config, meta *metadata.Metadata) error {
	wasmFilePath := filepath.Join(config.OutputDir, config.WasmFileName)

	bytes, err := wasmextractor.ReadWasmFile(wasmFilePath)
	if err != nil {
		return err
	}

	info, err := wasmextractor.ExtractWasmInfo(bytes)
	if err != nil {
		return err
	}

	// Remove unused imports (can easily happen when user code doesn't use all imports from a package)
	imports := make(map[string]bool, len(info.Imports))
	for _, i := range info.Imports {
		imports[i.Name] = true
	}
	for name := range meta.FnImports {
		if _, ok := imports[name]; !ok {
			delete(meta.FnImports, name)
		}
	}

	// Remove unused exports (less likely to happen, but still check)
	exports := make(map[string]bool, len(info.Exports))
	for _, e := range info.Exports {
		exports[e.Name] = true
	}
	for name := range meta.FnExports {
		if _, ok := exports[name]; !ok {
			delete(meta.FnExports, name)
		}
	}

	// Remove unused types (they might not be needed now, due to removing functions)
	var keptTypes = make(metadata.TypeMap, len(meta.Types))
	for _, fn := range append(utils.MapValues(meta.FnImports), utils.MapValues(meta.FnExports)...) {
		for _, param := range fn.Parameters {
			if _, ok := meta.Types[param.Type]; ok {
				keptTypes[param.Type] = meta.Types[param.Type]
				delete(meta.Types, param.Type)
			}
		}
		for _, result := range fn.Results {
			if _, ok := meta.Types[result.Type]; ok {
				keptTypes[result.Type] = meta.Types[result.Type]
				delete(meta.Types, result.Type)
			}
		}
	}

	// ensure types used by kept types are also kept
	for dirty := true; len(meta.Types) > 0 && dirty; {
		dirty = false

		keep := func(t string) {
			if _, ok := meta.Types[t]; ok {
				if _, ok := keptTypes[t]; !ok {
					keptTypes[t] = meta.Types[t]
					delete(meta.Types, t)
					dirty = true
				}
			}
		}

		for _, t := range keptTypes {
			if utils.IsPointerType(t.Name) {
				keep(utils.GetUnderlyingType(t.Name))
			} else if utils.IsListType(t.Name) {
				keep(utils.GetArraySubtype(t.Name))
			} else if utils.IsMapType(t.Name) {
				kt, vt := utils.GetMapSubtypes(t.Name)
				keep(kt)
				keep(vt)
			}

			for _, field := range t.Fields {
				keep(field.Type)
			}
		}
	}
	meta.Types = keptTypes

	return nil
}
