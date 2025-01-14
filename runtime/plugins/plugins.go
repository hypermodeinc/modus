/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package plugins

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/languages"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
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

	importsMap := make(map[string]wasm.FunctionDefinition, len(imports))
	for _, fnDef := range imports {
		if modName, fnName, ok := fnDef.Import(); ok {
			importName := modName + "." + fnName
			importsMap[importName] = fnDef
		}
	}

	for importName, fnMeta := range md.FnImports {
		fnDef, ok := importsMap[importName]
		if !ok {
			logger.Warn(ctx).Msgf("Unused import %s in plugin metadata. Please update your Modus SDK.", importName)
			continue
		}

		plan, err := planner.GetPlan(ctx, fnMeta, fnDef)
		if err != nil {
			return nil, fmt.Errorf("failed to get execution plan for %s: %w", importName, err)
		}
		plans[importName] = plan
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
