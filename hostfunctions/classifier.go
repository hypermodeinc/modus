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

type classifierResult struct {
	Label         string            `json:"label"`
	Confidence    float64           `json:"confidence"`
	Probabilities []classifierLabel `json:"probabilities"`
}

type classifierLabel struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
}

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	mem := mod.Memory()

	modelName, err := assemblyscript.ReadString(mem, pModelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading model name from wasm memory.")
		return 0
	}

	sentenceMapStr, err := assemblyscript.ReadString(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading sentence map string from wasm memory.")
		return 0
	}

	model, err := models.GetModel(modelName, manifest.GenerationTask)
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

	sentenceMap := make(map[string]string)
	if err := json.Unmarshal([]byte(sentenceMapStr), &sentenceMap); err != nil {
		logger.Err(ctx, err).Msg("Error unmarshalling sentence map.")
		return 0
	}

	result, err := models.PostToModelEndpoint[classifierResult](ctx, sentenceMap, model, host)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		logger.Err(ctx, err).Msg("Empty result returned from model.")
		return 0
	}

	resultMap := make(map[string]map[string]float64)
	for k, v := range result {
		resultMap[k] = make(map[string]float64)
		for _, p := range v.Probabilities {
			resultMap[k][p.Label] = p.Probability
		}
	}

	resBytes, err := json.Marshal(resultMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling classification result.")
		return 0
	}

	offset, err := assemblyscript.WriteString(ctx, mod, string(resBytes))
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
