/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"

	"hypruntime/plugins/metadata"
	"hypruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func (wa *wasmAdapter) InvokeFunction(ctx context.Context, function *metadata.Function, parameters map[string]any) (result any, err error) {

	// Get the wasm function
	fn := wa.mod.ExportedFunction(function.Name)
	if fn == nil {
		return nil, fmt.Errorf("function %s not found in wasm module", function.Name)
	}

	// Get parameters to pass as input to the function
	def := fn.Definition()
	params, indirect, cln, err := wa.getParameters(ctx, def, function, parameters)
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

	// Call the function
	res, err := fn.Call(ctx, params...)
	if err != nil {
		return nil, err
	}

	switch len(function.Results) {
	case 0:
		// no results are expected
		return nil, nil
	case 1:
		// a single result is expected
		resultType := function.Results[0].Type
		if len(res) == 1 {
			return wa.decodeObject(ctx, resultType, res)
		} else if indirect {
			return wa.readObject(ctx, resultType, uint32(params[0]))
		} else {
			// no actual result value, but we need to return a zero value of the expected type
			return wa.decodeObject(ctx, resultType, []uint64{0})
		}
	}

	// multiple results are expected (indirect)
	return wa.ReadIndirectResults(ctx, function, uint32(params[0]))
}

func (wa *wasmAdapter) getParameters(ctx context.Context, def wasm.FunctionDefinition, fn *metadata.Function, parameters map[string]any) ([]uint64, bool, utils.Cleaner, error) {

	params := make([]uint64, 0, len(def.ParamTypes()))
	cleaner := utils.NewCleanerN(len(fn.Parameters) + 1)
	indirect := false

	// If we expect results but the function signature doesn't have any, then TinyGo expects to be passed
	// a pointer in the first parameter, which indicates where the results should be stored.
	if len(def.ResultTypes()) == 0 && len(fn.Results) > 0 {

		totalSize := uint32(0)
		for _, r := range fn.Results {
			size, err := wa.typeInfo.GetSizeOfType(ctx, r.Type)
			if err != nil {
				return nil, false, nil, err
			}
			totalSize += size
		}

		// If totalSize is zero, then we have an edge case where there is no result value.
		// For example, a function that returns a struct with no fields, or a zero-length array.

		if totalSize > 0 {
			ptr, cln, err := wa.allocateWasmMemory(ctx, totalSize)
			if err != nil {
				return nil, true, cln, fmt.Errorf("failed to allocate memory for results: %w", err)
			}

			cleaner.AddCleaner(cln)
			params = append(params, uint64(ptr))
			indirect = true
		}
	}

	// Now we can encode the actual parameters
	for _, p := range fn.Parameters {
		val, found := parameters[p.Name]
		if !found && p.Default != nil {
			val = *p.Default
		}

		if val == nil && !wa.typeInfo.IsNullable(p.Type) {
			return nil, indirect, cleaner, fmt.Errorf("parameter '%s' cannot be null", p.Name)
		}

		encVals, cln, err := wa.encodeObject(ctx, p.Type, val)
		cleaner.AddCleaner(cln)
		if err != nil {
			return nil, indirect, cleaner, fmt.Errorf("function parameter '%s' is invalid: %w", p.Name, err)
		}

		params = append(params, encVals...)
	}

	return params, indirect, cleaner, nil
}

func (wa *wasmAdapter) ReadIndirectResults(ctx context.Context, fn *metadata.Function, offset uint32) (map[string]any, error) {
	results := make(map[string]any, len(fn.Results))

	fieldOffset := uint32(0)
	for i, r := range fn.Results {
		size, err := wa.typeInfo.GetSizeOfType(ctx, r.Type)
		if err != nil {
			return nil, err
		}

		fieldOffset += wa.typeInfo.getAlignmentPadding(fieldOffset, size)
		val, err := wa.readObject(ctx, r.Type, offset+fieldOffset)
		if err != nil {
			return nil, err
		}

		name := r.Name
		if name == "" {
			name = fmt.Sprintf("item%d", i+1)
		}

		results[name] = val
		fieldOffset += size
	}

	return results, nil
}

func (wa *wasmAdapter) WriteIndirectResults(ctx context.Context, fn *metadata.Function, offset uint32, results []any) (err error) {

	cleaner := utils.NewCleanerN(len(results))
	defer func() {
		if e := cleaner.Clean(); e != nil && err == nil {
			err = e
		}
	}()

	fieldOffset := uint32(0)
	for i, r := range results {
		typ := fn.Results[i].Type
		size, err := wa.typeInfo.GetSizeOfType(ctx, typ)
		if err != nil {
			return err
		}

		fieldOffset += wa.typeInfo.getAlignmentPadding(fieldOffset, size)
		cln, err := wa.writeObject(ctx, typ, offset+fieldOffset, r)
		cleaner.AddCleaner(cln)
		if err != nil {
			return err
		}

		fieldOffset += size
	}

	return nil
}
