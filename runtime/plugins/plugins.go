/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/runtime/hostfunctions/compatibility"
	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/languages"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/plugins/metadata"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Id             string
	Module         wazero.CompiledModule
	Metadata       *metadata.Metadata
	FileName       string
	Language       langsupport.Language
	ExecutionPlans map[string]langsupport.ExecutionPlan
}

func NewPlugin(ctx context.Context, cm wazero.CompiledModule, filename string, md *metadata.Metadata) (*Plugin, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	language, err := languages.GetLanguageForSDK(md.SDK)
	if err != nil {
		return nil, err
	}

	planner := language.NewPlanner(md)
	imports := cm.ImportedFunctions()
	exports := cm.ExportedFunctions()
	plans := make(map[string]langsupport.ExecutionPlan, len(imports)+len(exports))

	ctx = context.WithValue(ctx, utils.MetadataContextKey, md)

	for fnName, fnMeta := range md.FnExports {
		fnDef, ok := exports[fnName]
		if !ok {
			return nil, fmt.Errorf("no wasm function definition found for %s", fnName)
		}

		plan, err := planner.GetPlan(ctx, fnMeta, fnDef)
		if err != nil {
			return nil, fmt.Errorf("failed to get execution plan for %s: %w", fnName, err)
		}
		plans[fnName] = plan
	}

	warn := true
	for _, fnDef := range imports {
		modName, fnName, _ := fnDef.Import()
		impName := fmt.Sprintf("%s.%s", modName, fnName)

		fnMeta, err := getImportMetadata(ctx, modName, fnName, md, &warn)
		if err != nil {
			return nil, err
		} else if fnMeta == nil {
			continue
		}

		plan, err := planner.GetPlan(ctx, fnMeta, fnDef)
		if err != nil {
			return nil, fmt.Errorf("failed to get execution plan for %s: %w", impName, err)
		}
		plans[impName] = plan
	}

	plugin := &Plugin{
		Id:             utils.GenerateUUIDv7(),
		Module:         cm,
		Metadata:       md,
		FileName:       filename,
		Language:       language,
		ExecutionPlans: plans,
	}

	return plugin, nil
}

func getImportMetadata(ctx context.Context, modName, fnName string, md *metadata.Metadata, warn *bool) (*metadata.Function, error) {
	impName := fmt.Sprintf("%s.%s", modName, fnName)
	if fnMeta, ok := md.FnImports[impName]; ok {
		return fnMeta, nil
	}

	if modName == "hypermode" {
		if *warn {
			*warn = false
			logger.Warn(ctx).
				Str("plugin", md.Name()).
				Str("build_id", md.BuildId).
				Msg("Function imports are missing from the metadata. Using compatibility shims. Please update your SDK to the latest version.")
		}
		if fnMeta, err := compatibility.GetImportMetadataShim(impName); err != nil {
			return nil, fmt.Errorf("error creating compatibility shim for %s: %w", impName, err)
		} else {
			return fnMeta, nil
		}
	} else if shouldIgnoreModule(modName) {
		return nil, nil
	}

	return nil, fmt.Errorf("no metadata found for import %s", impName)
}

func shouldIgnoreModule(name string) bool {
	switch name {
	case "wasi_snapshot_preview1", "wasi", "env", "runtime", "syscall", "test":
		return true
	}

	return strings.HasPrefix(name, "wasi_")
}

func (p *Plugin) NameAndVersion() (name string, version string) {
	return p.Metadata.NameAndVersion()
}

func (p *Plugin) Name() string {
	return p.Metadata.Name()
}

func (p *Plugin) Version() string {
	return p.Metadata.Version()
}

func (p *Plugin) BuildId() string {
	return p.Metadata.BuildId
}

func GetPluginFromContext(ctx context.Context) (*Plugin, bool) {
	p, ok := ctx.Value(utils.PluginContextKey).(*Plugin)
	return p, ok
}
