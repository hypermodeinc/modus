/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

import (
	"context"
	"fmt"
	"runtime/debug"

	"hypruntime/plugins/metadata"
	"hypruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

type ExecutionPlan interface {
	FnDefinition() wasm.FunctionDefinition
	FnMetadata() *metadata.Function
	ParamHandlers() []TypeHandler
	ResultHandlers() []TypeHandler
	UseResultIndirection() bool
	HasDefaultParameters() bool
	InvokeFunction(ctx context.Context, wa WasmAdapter, parameters map[string]any) (result any, err error)
}

func NewExecutionPlan(fnDef wasm.FunctionDefinition, fnMeta *metadata.Function, paramHandlers, resultHandlers []TypeHandler, indirectResultSize uint32) ExecutionPlan {
	hasDefaultParameters := false
	for _, p := range fnMeta.Parameters {
		if p.Default != nil {
			hasDefaultParameters = true
			break
		}
	}

	return &executionPlan{fnDef, fnMeta, paramHandlers, resultHandlers, indirectResultSize, hasDefaultParameters}
}

type executionPlan struct {
	fnDefinition         wasm.FunctionDefinition
	fnMetadata           *metadata.Function
	paramHandlers        []TypeHandler
	resultHandlers       []TypeHandler
	indirectResultSize   uint32
	hasDefaultParameters bool
}

func (p *executionPlan) FnDefinition() wasm.FunctionDefinition {
	return p.fnDefinition
}

func (p *executionPlan) FnMetadata() *metadata.Function {
	return p.fnMetadata
}

func (p *executionPlan) ParamHandlers() []TypeHandler {
	return p.paramHandlers
}

func (p *executionPlan) ResultHandlers() []TypeHandler {
	return p.resultHandlers
}

func (p *executionPlan) UseResultIndirection() bool {
	return p.indirectResultSize > 0
}

func (p *executionPlan) HasDefaultParameters() bool {
	return p.hasDefaultParameters
}

func (plan *executionPlan) InvokeFunction(ctx context.Context, wa WasmAdapter, parameters map[string]any) (result any, err error) {
	// Recover from panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			err = utils.ConvertToError(r)
			if utils.HypermodeDebugEnabled() {
				debug.PrintStack()
			}
		}
	}()

	// Get the wasm function
	fnName := plan.FnMetadata().Name
	fn := wa.GetFunction(fnName)
	if fn == nil {
		return nil, fmt.Errorf("function %s not found in wasm module", fnName)
	}

	// Get parameters to pass as input to the function
	params, cln, err := plan.getWasmParameters(ctx, wa, parameters)
	defer func() {
		// Clean up any resources allocated for the parameters (when done)
		if cln != nil {
			if e := cln.Clean(); e != nil && err == nil {
				err = e
			}
		}
	}()
	if err != nil {
		return nil, err
	}

	// Pre-invoke hook
	if err := wa.PreInvoke(ctx, plan); err != nil {
		return nil, err
	}

	// Call the function
	res, err := fn.Call(ctx, params...)
	if err != nil {
		return nil, err
	}

	// Get the result indirection pointer (if any)
	var indirectPtr uint32
	if plan.UseResultIndirection() {
		indirectPtr = uint32(params[0])
	}

	// Interpret and return the results
	return plan.interpretWasmResults(ctx, wa, res, indirectPtr)
}

func (plan *executionPlan) getWasmParameters(ctx context.Context, wa WasmAdapter, parameters map[string]any) ([]uint64, utils.Cleaner, error) {
	var paramVals []uint64
	var cleaner utils.Cleaner

	if plan.UseResultIndirection() {
		ptr, cln, err := wa.AllocateMemory(ctx, plan.indirectResultSize)
		if err != nil {
			return nil, cln, fmt.Errorf("failed to allocate memory for results: %w", err)
		}

		paramVals = make([]uint64, 1, len(plan.FnDefinition().ParamTypes())+1)
		cleaner = utils.NewCleanerN(len(plan.ParamHandlers()) + 1)
		paramVals[0] = uint64(ptr)
		cleaner.AddCleaner(cln)
	} else {
		paramVals = make([]uint64, 0, len(plan.FnDefinition().ParamTypes()))
		cleaner = utils.NewCleanerN(len(plan.ParamHandlers()))
	}

	var mask uint64
	var has_opt bool
	handlers := plan.ParamHandlers()
	for i, p := range plan.FnMetadata().Parameters {

		val, found := parameters[p.Name]
		if found {
			mask |= 1 << i
		} else if p.Default != nil {
			val = *p.Default
			mask |= 1 << i
		} else if p.Optional {
			has_opt = true
			continue
		}

		encVals, cln, err := handlers[i].Encode(ctx, wa, val)
		cleaner.AddCleaner(cln)
		if err != nil {
			return nil, cleaner, fmt.Errorf("function parameter '%s' is invalid: %w", p.Name, err)
		}

		paramVals = append(paramVals, encVals...)
	}

	if has_opt {
		paramVals = append(paramVals, mask)
	}

	return paramVals, cleaner, nil
}

func (plan *executionPlan) interpretWasmResults(ctx context.Context, wa WasmAdapter, vals []uint64, indirectPtr uint32) (any, error) {

	handlers := plan.ResultHandlers()
	switch len(handlers) {
	case 0:
		// no results are expected
		return nil, nil
	case 1:
		// a single result is expected
		handler := handlers[0]
		if plan.UseResultIndirection() {
			return handler.Read(ctx, wa, indirectPtr)
		} else if len(vals) == 1 {
			return handler.Decode(ctx, wa, vals)
		} else {
			// no actual result value, but we need to return a zero value of the expected type
			return handler.TypeInfo().ZeroValue(), nil
		}
	}

	// multiple results are expected (indirect)
	return plan.readIndirectResults(ctx, wa, indirectPtr)
}

func (plan *executionPlan) readIndirectResults(ctx context.Context, wa WasmAdapter, offset uint32) ([]any, error) {

	// multiple-results are read like a struct

	handlers := plan.ResultHandlers()
	results := make([]any, len(handlers))

	fieldOffset := uint32(0)
	for i, handler := range handlers {
		size := handler.TypeInfo().Size()

		fieldType := handler.TypeInfo().Name()
		alignment, err := wa.TypeInfo().GetAlignmentOfType(ctx, fieldType)
		if err != nil {
			return nil, err
		}

		fieldOffset = AlignOffset(fieldOffset, alignment)

		val, err := handler.Read(ctx, wa, offset+fieldOffset)
		if err != nil {
			return nil, err
		}

		results[i] = val
		fieldOffset += size
	}

	return results, nil
}
