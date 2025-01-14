/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package extractor

import (
	"go/types"
	"sort"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/config"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metadata"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
)

func CollectProgramInfo(config *config.Config, meta *metadata.Metadata) error {
	pkgs, err := loadPackages(config.SourceDir)
	if err != nil {
		return err
	}

	requiredTypes := make(map[string]types.Type)

	for name, f := range getExportedFunctions(pkgs) {
		meta.FnExports[name] = transformFunc(name, f, pkgs)
		findRequiredTypes(f, requiredTypes)
	}

	for name, f := range getImportedFunctions(pkgs) {
		meta.FnImports[name] = transformFunc(name, f, pkgs)
		findRequiredTypes(f, requiredTypes)
	}

	// proxy imports overwrite regular imports
	for name, f := range getProxyImportFunctions(pkgs) {
		if _, ok := meta.FnImports[name]; ok {
			meta.FnImports[name] = transformFunc(name, f, pkgs)
			findRequiredTypes(f, requiredTypes)
		}
	}

	id := uint32(3) // 1 and 2 are reserved for []byte and string
	keys := utils.MapKeys(requiredTypes)
	sort.Strings(keys)
	for _, name := range keys {
		t := requiredTypes[name]
		if s, ok := t.(*types.Struct); ok && !wellKnownTypes[name] {
			t := transformStruct(name, s, pkgs)
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
