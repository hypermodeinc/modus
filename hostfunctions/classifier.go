/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/manifest"
	"hmruntime/models"

	wasm "github.com/tetratelabs/wazero/api"
)

type classifierResult struct {
	Label         string            `json:"label"`
	Confidence    float32           `json:"confidence"`
	Probabilities []classifierLabel `json:"probabilities"`
}

type classifierLabel struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {

	var modelName string
	var sentenceMap map[string]string
	err := readParams2(ctx, mod, pModelName, pSentenceMap, &modelName, &sentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	model, err := models.GetModel(modelName, manifest.ClassificationTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	result, err := models.PostToModelEndpoint[classifierResult](ctx, sentenceMap, model)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		logger.Err(ctx, err).Msg("Empty result returned from model.")
		return 0
	}

	resultMap := make(map[string]map[string]float32)
	for k, v := range result {
		resultMap[k] = make(map[string]float32)
		for _, p := range v.Probabilities {
			resultMap[k][p.Label] = p.Probability
		}
	}

	offset, err := writeResult(ctx, mod, resultMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing classification result.")
		return 0
	}

	return offset
}
