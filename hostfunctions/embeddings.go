/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"encoding/json"

	"hmruntime/functions/assemblyscript"
	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/manifest"
	"hmruntime/models"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, manifest.EmbeddingTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	var host manifest.Host
	if model.Host != hypermodeHostName {
		host, err = hosts.GetHost(model.Host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting model host.")
			return 0
		}
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}

	result, err := models.PostToModelEndpoint[[]float64](ctx, sentenceMap, model, host)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		logger.Error(ctx).Msg("Empty result returned from model.")
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling embedding result.")
		return 0
	}

	offset, err := assemblyscript.WriteString(ctx, mod, string(res))
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
