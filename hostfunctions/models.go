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

func hostLookupModel(ctx context.Context, mod wasm.Module, pModelName uint32) (pModelInfo uint32) {

	// Read input parameters
	var modelName string
	err := readParams(ctx, mod, param{pModelName, &modelName})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
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
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, info)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

func hostInvokeModel(ctx context.Context, mod wasm.Module, pModelName, pInput uint32) (pOutput uint32) {

	// Read input parameters
	var modelName, input string
	err := readParams(ctx, mod, param{pModelName, &modelName}, param{pInput, &input})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
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
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
