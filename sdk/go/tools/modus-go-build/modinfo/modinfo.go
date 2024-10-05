/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package modinfo

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/config"

	"github.com/hashicorp/go-version"
	"golang.org/x/mod/modfile"
)

const sdkModulePath = "github.com/hypermodeinc/modus/sdk/go"

type ModuleInfo struct {
	ModulePath      string
	ModusSDKVersion *version.Version
}

func CollectModuleInfo(config *config.Config) (*ModuleInfo, error) {
	modFilePath := filepath.Join(config.SourceDir, "go.mod")
	data, err := os.ReadFile(modFilePath)
	if err != nil {
		return nil, err
	}

	modFile, err := modfile.Parse(modFilePath, data, nil)
	if err != nil {
		return nil, err
	}

	modPath := modFile.Module.Mod.Path
	config.WasmFileName = path.Base(modPath) + ".wasm"

	result := &ModuleInfo{
		ModulePath: modPath,
	}

	for _, requiredModules := range modFile.Require {
		mod := requiredModules.Mod
		if mod.Path == sdkModulePath {
			sVer := strings.TrimPrefix(mod.Version, "v")
			if ver, err := version.NewVersion(sVer); err == nil {
				result.ModusSDKVersion = ver
			}
		}
	}

	return result, nil
}
