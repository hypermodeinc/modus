/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
)

func PostToModelEndpoint_Old[TResult any](ctx context.Context, sentenceMap map[string]string, model manifest.ModelInfo) (map[string]TResult, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
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

	res, err := postToModelEndpoint[PredictionResult[TResult]](ctx, model, req)
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

type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJson OutputFormat = "json_object"
)

type PredictionResult[T any] struct {
	Predictions []T `json:"predictions"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []MessageChoice `json:"choices"`
	Error   InvokeError     `json:"error"`
}

type MessageChoice struct {
	Message ChatMessage `json:"message"`
}

type InvokeError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}
