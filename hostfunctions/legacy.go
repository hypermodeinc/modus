/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/models/legacymodels"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	var modelName string
	var sentenceMap map[string]string
	err := readParams2(ctx, mod, pModelName, pSentenceMap, &modelName, &sentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	resultMap, err := legacymodels.InvokeClassifier(ctx, modelName, sentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error invoking classifier.")
		return 0
	}

	offset, err := writeResult(ctx, mod, resultMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing classification result.")
		return 0
	}

	return offset
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {

	var modelName string
	var sentenceMap map[string]string
	err := readParams2(ctx, mod, pModelName, pSentenceMap, &modelName, &sentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	result, err := legacymodels.ComputeEmbedding(ctx, modelName, sentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error invoking classifier.")
		return 0
	}

	offset, err := writeResult(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing classification result.")
		return 0
	}

	return offset
}

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, pModelName uint32, pInstruction uint32, pSentence uint32, pFormat uint32) uint32 {

	var modelName, instruction, sentence, format string
	err := readParams4(ctx, mod, pModelName, pInstruction, pSentence, pFormat, &modelName, &instruction, &sentence, &format)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	content, err := legacymodels.InvokeTextGenerator(ctx, modelName, format, instruction, sentence)
	if err != nil {
		logger.Error(ctx).Msg("Error invoking text generator.")
		return 0
	}

	offset, err := writeResult(ctx, mod, content)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
