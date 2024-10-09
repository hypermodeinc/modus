/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package legacymodels

import (
	"context"
	"errors"

	"github.com/hypermodeinc/modus/runtime/models"
)

func InvokeClassifier(ctx context.Context, modelName string, sentenceMap map[string]string) (map[string]map[string]float32, error) {
	model, err := models.GetModel(modelName)
	if err != nil {
		return nil, err
	}

	result, err := postToModelEndpoint[classifierResult](ctx, model, sentenceMap)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("empty result returned from model")
	}

	resultMap := make(map[string]map[string]float32)
	for k, v := range result {
		resultMap[k] = make(map[string]float32)
		for _, p := range v.Probabilities {
			resultMap[k][p.Label] = p.Probability
		}
	}

	return resultMap, nil
}

type classifierResult struct {
	Label         string            `json:"label"`
	Confidence    float32           `json:"confidence"`
	Probabilities []classifierLabel `json:"probabilities"`
}

type classifierLabel struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}
