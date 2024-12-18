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

	return nil
}
