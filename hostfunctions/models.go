/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/logger"
	"hmruntime/models"

	wasm "github.com/tetratelabs/wazero/api"
)

func init() {
	addHostFunction(&hostFunctionDefinition{
		name:     "lookupModel",
		function: wasm.GoModuleFunc(hostLookupModel),
		params:   []wasm.ValueType{wasm.ValueTypeI32},
		results:  []wasm.ValueType{wasm.ValueTypeI32},
	})

	addHostFunction(&hostFunctionDefinition{
		name:     "invokeModel",
		function: wasm.GoModuleFunc(hostInvokeModel),
		params:   []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		results:  []wasm.ValueType{wasm.ValueTypeI32},
	})
}

func hostLookupModel(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var modelName string
	if err := readParams(ctx, mod, stack, &modelName); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled getting model info.",
		Error:     "Error getting model info.",
		Detail:    fmt.Sprintf("Model: %s", modelName),
	}

	// Prepare the host function
	var info *models.ModelInfo
	fn := func() (err error) {
		info, err = models.GetModelInfo(modelName)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, info); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostInvokeModel(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var modelName, input string
	if err := readParams(ctx, mod, stack, &modelName, &input); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Invoking model.",
		Completed: "Completed model invocation.",
		Cancelled: "Cancelled model invocation.",
		Error:     "Error invoking model.",
		Detail:    fmt.Sprintf("Model: %s", modelName),
	}

	// Prepare the host function
	var output string
	fn := func() (err error) {
		output, err = models.InvokeModel(ctx, modelName, input)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, output); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}
