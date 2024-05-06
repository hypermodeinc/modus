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

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {

	var modelName string
	var sentenceMap map[string]string
	err := readParams2(ctx, mod, pModelName, pSentenceMap, &modelName, &sentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	var result map[string][]float64
	inferenceService, err := models.CreateInferenceService(modelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error instanciating inference service.")
		return 0
	}
	result, err = inferenceService.ComputeEmbedding(ctx, sentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error in ComputeEmbedding")
		return 0
	}

	if len(result) == 0 {
		logger.Error(ctx).Msg("Empty result returned from model.")
		return 0
	}

	offset, err := writeResult(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing classification result.")
		return 0
	}

	return offset
}
