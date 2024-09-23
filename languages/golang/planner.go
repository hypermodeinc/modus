/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"

	"hypruntime/langsupport"
	"hypruntime/plugins/metadata"

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
	} else if ti.IsPointer() {
		return p.NewPointerHandler(ctx, ti)
	} else if ti.IsList() {
		if _langTypeInfo.IsSliceType(typeName) {
			if ti.ListElementType().IsPrimitive() {
				return p.NewPrimitiveSliceHandler(ti)
			} else {
				return p.NewSliceHandler(ctx, ti)
			}
		} else if _langTypeInfo.IsArrayType(typeName) {
			if ti.ListElementType().IsPrimitive() {
				return p.NewPrimitiveArrayHandler(ti)
			} else {
				return p.NewArrayHandler(ctx, ti)
			}
		}
	} else if ti.IsMap() {
		return p.NewMapHandler(ctx, ti)
	} else if ti.IsTimestamp() {
		return p.NewTimeHandler(ti)
	} else if ti.IsObject() {
		return p.NewStructHandler(ctx, ti)
	}

	return nil, fmt.Errorf("can't determine plan for type: %s", typeName)
}

func (p *planner) GetPlan(ctx context.Context, fnMeta *metadata.Function, fnDef wasm.FunctionDefinition) (langsupport.ExecutionPlan, error) {

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

	indirectResultSize, err := p.getIndirectResultSize(ctx, fnMeta, fnDef)
	if err != nil {
		return nil, err
	}

	plan := langsupport.NewExecutionPlan(fnDef, fnMeta, paramHandlers, resultHandlers, indirectResultSize)
	return plan, nil
}

func (p *planner) getIndirectResultSize(ctx context.Context, fnMeta *metadata.Function, fnDef wasm.FunctionDefinition) (uint32, error) {

	// If no results are expected, then we don't need to use indirection.
	if len(fnMeta.Results) == 0 {
		return 0, nil
	}

	// If the function definition has results, then we don't need to use indirection.
	if len(fnDef.ResultTypes()) > 0 {
		return 0, nil
	}

	// We expect results but the function signature doesn't have any.
	// Thus, TinyGo expects to be passed a pointer in the first parameter,
	// which indicates where the results should be stored.
	//
	// However, if totalSize is zero, then we have an edge case where there is no result value.
	// For example, a function that returns a struct with no fields, or a zero-length array.
	//
	// We need the total size either way, because we will need to allocate memory for the results.

	totalSize := uint32(0)
	for _, r := range fnMeta.Results {
		size, err := _langTypeInfo.GetSizeOfType(ctx, r.Type)
		if err != nil {
			return 0, err
		}
		totalSize += size
	}
	return totalSize, nil
}
