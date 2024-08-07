/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/logger"
	"hmruntime/models/legacymodels"

	wasm "github.com/tetratelabs/wazero/api"
)

func init() {
	addHostFunction(&hostFunctionDefinition{
		name:     "invokeClassifier",
		function: wasm.GoModuleFunc(hostInvokeClassifier),
		params:   []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		results:  []wasm.ValueType{wasm.ValueTypeI32},
	})

	addHostFunction(&hostFunctionDefinition{
		name:     "computeEmbedding",
		function: wasm.GoModuleFunc(hostComputeEmbedding),
		params:   []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		results:  []wasm.ValueType{wasm.ValueTypeI32},
	})

	addHostFunction(&hostFunctionDefinition{
		name:     "invokeTextGenerator",
		function: wasm.GoModuleFunc(hostInvokeTextGenerator),
		params:   []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		results:  []wasm.ValueType{wasm.ValueTypeI32},
	})
}

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var modelName string
	var sentenceMap map[string]string
	if err := readParams(ctx, mod, stack, &modelName, &sentenceMap); err != nil {
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
	var resultMap map[string]map[string]float32
	fn := func() (err error) {
		resultMap, err = legacymodels.InvokeClassifier(ctx, modelName, sentenceMap)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, resultMap); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var modelName string
	var sentenceMap map[string]string
	if err := readParams(ctx, mod, stack, &modelName, &sentenceMap); err != nil {
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
	var result map[string][]float64
	fn := func() (err error) {
		result, err = legacymodels.ComputeEmbedding(ctx, modelName, sentenceMap)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, result); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var modelName, instruction, sentence, format string
	if err := readParams(ctx, mod, stack, &modelName, &instruction, &sentence, &format); err != nil {
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
	var content string
	fn := func() (err error) {
		content, err = legacymodels.InvokeTextGenerator(ctx, modelName, format, instruction, sentence)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, content); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}
