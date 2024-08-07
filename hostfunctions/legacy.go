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

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, pModelName, pSentenceMap uint32) uint32 {

	// Read input parameters
	var modelName string
	var sentenceMap map[string]string
	err := readParams(ctx, mod, param{pModelName, &modelName}, param{pSentenceMap, &sentenceMap})
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
	var resultMap map[string]map[string]float32
	fn := func() (err error) {
		resultMap, err = legacymodels.InvokeClassifier(ctx, modelName, sentenceMap)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, resultMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName, pSentenceMap uint32) uint32 {

	// Read input parameters
	var modelName string
	var sentenceMap map[string]string
	err := readParams(ctx, mod, param{pModelName, &modelName}, param{pSentenceMap, &sentenceMap})
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
	var result map[string][]float64
	fn := func() (err error) {
		result, err = legacymodels.ComputeEmbedding(ctx, modelName, sentenceMap)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, pModelName, pInstruction, pSentence, pFormat uint32) uint32 {

	// Read input parameters
	var modelName, instruction, sentence, format string
	err := readParams(ctx, mod, param{pModelName, &modelName}, param{pInstruction, &instruction}, param{pSentence, &sentence}, param{pFormat, &format})
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
	var content string
	fn := func() (err error) {
		content, err = legacymodels.InvokeTextGenerator(ctx, modelName, format, instruction, sentence)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, content)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
