/*
 * Copyright 2024 Hypermode, Inc.
 */

package legacymodels

import (
	"context"
	"errors"

	"hmruntime/models"
)

func ComputeEmbedding(ctx context.Context, modelName string, sentenceMap map[string]string) (map[string][]float64, error) {
	model, err := models.GetModel(modelName)
	if err != nil {
		return nil, err
	}

	result, err := postToModelEndpoint[[]float64](ctx, model, sentenceMap)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("empty result returned from model")
	}

	return result, nil
}
