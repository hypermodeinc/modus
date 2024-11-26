/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package metagen

import (
	"fmt"
	"path"

	"os"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/config"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/extractor"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/gitinfo"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metadata"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/modinfo"
)

const sdkName = "modus-sdk-go"

func GenerateMetadata(config *config.Config, mod *modinfo.ModuleInfo) (*metadata.Metadata, error) {
	if _, err := os.Stat(config.SourceDir); err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	meta := metadata.NewMetadata()
	meta.Module = mod.ModulePath
	meta.Plugin = path.Base(mod.ModulePath)

	meta.SDK = sdkName
	if mod.ModusSDKVersion != nil {
		meta.SDK += "@" + mod.ModusSDKVersion.String()
	}

	if err := extractor.CollectProgramInfo(config, meta); err != nil {
		return nil, fmt.Errorf("error collecting program info: %w", err)
	}

	gitinfo.TryCollectGitInfo(config, meta)

	return meta, nil
}
