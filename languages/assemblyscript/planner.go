/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"

	"hypruntime/langsupport"
	"hypruntime/plugins/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

func NewPlanner(metadata *metadata.Metadata) langsupport.Planner {
	return &planner{
		typeHandlers: make(map[string]langsupport.TypeHandler),
		metadata:     metadata,
	}
}

type planner struct {
	typeHandlers map[string]langsupport.TypeHandler
	metadata     *metadata.Metadata
}

func NewTypeHandler[T any](p *planner, typ string) *T {
	handler := new(T)
	p.typeHandlers[typ] = any(handler).(langsupport.TypeHandler)
	return handler
}

func (p *planner) AllHandlers() map[string]langsupport.TypeHandler {
	return p.typeHandlers
}

func (p *planner) GetHandler(ctx context.Context, typ string) (langsupport.TypeHandler, error) {
	if handler, ok := p.typeHandlers[typ]; ok {
		return handler, nil
	}

	rt, err := _typeInfo.GetReflectedType(ctx, typ)
	if err != nil {
		return nil, err
	}

	if _typeInfo.IsPrimitiveType(typ) {
		return p.NewPrimitiveHandler(typ, rt)
	} else if _typeInfo.IsStringType(typ) {
		return p.NewStringHandler(typ, rt)
	} else if _typeInfo.IsArrayBufferType(typ) {
		return p.NewArrayBufferHandler(typ, rt)
	} else {
		return p.NewManagedObjectHandler(ctx, typ, rt)
	}
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

	plan := langsupport.NewExecutionPlan(fnDef, fnMeta, paramHandlers, resultHandlers, 0)
	return plan, nil
}
