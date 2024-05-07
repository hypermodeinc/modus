/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/models"

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
	// Get the inference service for the model and invoke the classifier function
	inferenceService, err := models.CreateInferenceService(modelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error instanciating inference service.")
		return 0
	}
	resultMap, err := inferenceService.InvokeClassifier(ctx, sentenceMap)

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
