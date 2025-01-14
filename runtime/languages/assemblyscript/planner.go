/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func NewPlanner(metadata *metadata.Metadata) langsupport.Planner {
	return &planner{
		typeCache:    make(map[string]langsupport.TypeInfo),
		typeHandlers: make(map[string]langsupport.TypeHandler),
		metadata:     metadata,
	}
}

type planner struct {
	typeCache    map[string]langsupport.TypeInfo
	typeHandlers map[string]langsupport.TypeHandler
	metadata     *metadata.Metadata
}

func (p *planner) AddHandler(h langsupport.TypeHandler) {
	p.typeHandlers[h.TypeInfo().Name()] = h
}

func (p *planner) AllHandlers() map[string]langsupport.TypeHandler {
	return p.typeHandlers
}

func NewTypeHandler(ti langsupport.TypeInfo) *typeHandler {
	return &typeHandler{
		typeInfo: ti,
	}
}

type typeHandler struct {
	typeInfo langsupport.TypeInfo
}

func (h *typeHandler) TypeInfo() langsupport.TypeInfo {
	return h.typeInfo
}

func (p *planner) GetHandler(ctx context.Context, typeName string) (langsupport.TypeHandler, error) {
	if handler, ok := p.typeHandlers[typeName]; ok {
		return handler, nil
	}

	ti, err := GetTypeInfo(ctx, typeName, p.typeCache)
	if err != nil {
		return nil, fmt.Errorf("failed to get type info for %s: %w", typeName, err)
	}

	if ti.IsPrimitive() {
		return p.NewPrimitiveHandler(ti)
	} else if ti.IsString() {
		return p.NewStringHandler(ti)
	} else if _langTypeInfo.IsArrayBufferType(typeName) {
		return p.NewArrayBufferHandler(ti)
	} else {
		return p.NewManagedObjectHandler(ctx, ti)
	}
}

func (p *planner) GetPlan(ctx context.Context, fnMeta *metadata.Function, fnDef wasm.FunctionDefinition) (langsupport.ExecutionPlan, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	paramHandlers := make([]langsupport.TypeHandler, len(fnMeta.Parameters))
	for i, param := range fnMeta.Parameters {
		handler, err := p.GetHandler(ctx, param.Type)
		if err != nil {
			return nil, err
		}
		paramHandlers[i] = handler
	}

	resultHandlers := make([]langsupport.TypeHandler, len(fnMeta.Results))
	for i, result := range fnMeta.Results {
		handler, err := p.GetHandler(ctx, result.Type)
		if err != nil {
			return nil, err
		}
		resultHandlers[i] = handler
	}

	plan := langsupport.NewExecutionPlan(fnDef, fnMeta, paramHandlers, resultHandlers, 0)
	return plan, nil
}
