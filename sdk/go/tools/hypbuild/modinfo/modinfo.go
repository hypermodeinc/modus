/*
 * Copyright 2024 Hypermode, Inc.
 */

package modinfo

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hypermodeAI/functions-go/tools/hypbuild/config"

	"github.com/hashicorp/go-version"
	"golang.org/x/mod/modfile"
)

const sdkModulePath = "github.com/hypermodeAI/functions-go"

type ModuleInfo struct {
	ModulePath          string
	HypermodeSDKVersion *version.Version
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
				result.HypermodeSDKVersion = ver
			}
		}
	}

	return result, nil
}
