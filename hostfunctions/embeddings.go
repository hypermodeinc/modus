/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/manifest"
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

	model, err := models.GetModel(modelName, manifest.EmbeddingTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}
	var host manifest.Host
	if model.Host != hosts.HypermodeHost {
		host, err = hosts.GetHost(model.Host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting model host.")
			return 0
		}
	}
	var result map[string][]float64
	if (models.OpenAIHost == model.Host) || (models.MistralHost == model.Host) {
		// use LLM
		llm, err := models.CreateLlmService(model.Host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error instanciating LLM")
			return 0
		}

		result, err = llm.Embedding(ctx, sentenceMap, model, host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error embeddings with LLM.")
			return 0
		}

	} else {
		result, err = models.PostToModelEndpoint[[]float64](ctx, sentenceMap, model)
		if err != nil {
			logger.Err(ctx, err).Msg("Error posting to model endpoint.")
			return 0
		}
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
