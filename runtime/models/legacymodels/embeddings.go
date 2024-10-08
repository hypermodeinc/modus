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
