/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package legacymodels

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/pkg/manifest"
	"github.com/hypermodeinc/modus/runtime/models"
	"github.com/hypermodeinc/modus/runtime/utils"
)

type predictionResult[T any] struct {
	Predictions []T `json:"predictions"`
}

func postToModelEndpoint[TResult any](ctx context.Context, model *manifest.ModelInfo, sentenceMap map[string]string) (map[string]TResult, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	// self hosted models takes in array, can optimize for parallelizing later
	keys, sentences := []string{}, []string{}

	for k, v := range sentenceMap {
		// create a map of sentences to send to the model
		sentences = append(sentences, v)
		// create a list of keys to map the results back to the original sentences
		keys = append(keys, k)
	}
	// create a map of sentences to send to the model
	req := map[string][]string{"instances": sentences}

	res, err := models.PostToModelEndpoint[predictionResult[TResult]](ctx, model, req)
	if err != nil {
		return nil, err
	}

	if len(res.Predictions) != len(keys) {
		return nil, fmt.Errorf("number of predictions does not match number of sentences")
	}

	// map the results back to the original sentences
	result := make(map[string]TResult)
	for i, v := range res.Predictions {
		result[keys[i]] = v
	}

	return result, nil
}
