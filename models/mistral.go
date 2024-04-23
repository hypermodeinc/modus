/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"
	"hmruntime/hosts"
	"hmruntime/manifest"
	"hmruntime/utils"
)

type mistral struct{}

func (llm *mistral) Embedding(ctx context.Context, sentenceMap map[string]string, model manifest.Model, host manifest.Host) (map[string][]float64, error) {
	// Get the API key to use for this model
	key, err := hosts.GetHostKey(ctx, host)
	if err != nil {
		return nil, err
	}
	// Convert map to slice of values.
	values := []string{}
	keys := []string{}
	for key, value := range sentenceMap {
		values = append(values, value)
		keys = append(keys, key)
	}

	// build the request body following OpenAI API
	reqBody := EmbeddingRequest{
		Model:          model.SourceModel,
		Input:          values,
		EncodingFormat: "float",
	}

	// We ignore the model endpoint and use the OpenAI endpoint
	const endpoint = "https://api.mistral.ai/v1/embeddings"
	headers := map[string]string{
		"Authorization": "Bearer " + key,
	}

	result, err := utils.PostHttp[EmbeddingResponse](endpoint, reqBody, headers)

	if err != nil {
		return nil, fmt.Errorf("error posting to %s: %w", endpoint, err)
	}

	// Convert result to map .
	resultMap := make(map[string][]float64)
	for _, value := range result.Data {
		resultMap[keys[value.Index]] = value.Embedding
	}

	return resultMap, nil
}

func (llm *mistral) ChatCompletion(ctx context.Context, model manifest.Model, host manifest.Host, instruction string, sentence string, outputFormat OutputFormat) (ChatResponse, error) {

	// Get the API key to use for this model
	key, err := hosts.GetHostKey(ctx, host)
	if err != nil {
		return ChatResponse{}, err
	}

	// build the request body following OpenAI API
	reqBody := ChatContext{
		Model: model.SourceModel,
		ResponseFormat: ResponseFormat{
			Type: string(outputFormat),
		},
		Messages: []ChatMessage{
			{Role: "system", Content: instruction},
			{Role: "user", Content: sentence},
		},
	}

	// We ignore the model endpoint and use the OpenAI endpoint
	const endpoint = "https://api.mistral.ai/v1/chat/completions"
	headers := map[string]string{
		"Authorization": "Bearer " + key,
	}

	result, err := utils.PostHttp[ChatResponse](endpoint, reqBody, headers)

	if err != nil {
		return ChatResponse{}, fmt.Errorf("error posting to %s: %w", endpoint, err)
	}

	if result.Error.Message != "" {
		return ChatResponse{}, fmt.Errorf("error returned from %s: %s", endpoint, result.Error.Message)
	}

	return result, nil
}
