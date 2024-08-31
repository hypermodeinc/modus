/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"
	"reflect"

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

	switch rt.Kind() {

	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr, reflect.UnsafePointer:

		return p.NewPrimitiveHandler(typ, rt)

	case reflect.String:
		return p.NewStringHandler(typ, rt)

	case reflect.Pointer:
		return p.NewPointerHandler(ctx, typ, rt)

	case reflect.Slice:
		if _typeInfo.IsPrimitiveType(_typeInfo.GetListSubtype(typ)) {
			return p.NewPrimitiveSliceHandler(typ, rt)
		} else {
			return p.NewSliceHandler(ctx, typ, rt)
		}

	case reflect.Array:
		if _typeInfo.IsPrimitiveType(_typeInfo.GetListSubtype(typ)) {
			return p.NewPrimitiveArrayHandler(typ, rt)
		} else {
			return p.NewArrayHandler(ctx, typ, rt)
		}

	case reflect.Map:
		if _typeInfo.IsMapType(typ) {
			return p.NewMapHandler(ctx, typ, rt)
		} else {
			// This is a struct that is being passed as a map.
			return p.NewStructHandler(ctx, typ, rt)
		}

	case reflect.Struct:
		if _typeInfo.IsTimestampType(typ) {
			return p.NewTimeHandler(typ, rt)
		} else {
			return p.NewStructHandler(ctx, typ, rt)
		}
	}

	return nil, fmt.Errorf("can't determine plan for type: %s", typ)
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
		size, err := _typeInfo.GetSizeOfType(ctx, r.Type)
		if err != nil {
			return 0, err
		}
		totalSize += size
	}
	return totalSize, nil
}
