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

type openai struct {
	model manifest.Model
	host  manifest.Host
}

func (llm *openai) InvokeClassifier(ctx context.Context, input []string) (map[string]float64, error) {
	return nil, fmt.Errorf("invokeClassifier not implemented for openai model")
}

func (llm *openai) ComputeEmbedding(ctx context.Context, sentenceMap map[string]string) (map[string][]float64, error) {
	return nil, fmt.Errorf("Embedding not implemented for Openai")
}

func (llm *openai) ChatCompletion(ctx context.Context, instruction string, sentence string, outputFormat OutputFormat) (ChatResponse, error) {

	// Get the OpenAI API key to use for this model
	key, err := hosts.GetHostKey(ctx, llm.host)
	if err != nil {
		return ChatResponse{}, err
	}

	// build the request body following OpenAI API
	reqBody := ChatContext{
		Model: llm.model.SourceModel,
		ResponseFormat: ResponseFormat{
			Type: string(outputFormat),
		},
		Messages: []ChatMessage{
			{Role: "system", Content: instruction},
			{Role: "user", Content: sentence},
		},
	}

	// We ignore the model endpoint and use the OpenAI endpoint
	const endpoint = "https://api.openai.com/v1/chat/completions"
	headers := map[string]string{
		"Authorization": "Bearer " + key,
	}

	result, err := utils.PostHttp[ChatResponse](endpoint, reqBody, headers)

	if err != nil {
		return ChatResponse{}, fmt.Errorf("error posting to OpenAI: %w", err)
	}

	if result.Error.Message != "" {
		return ChatResponse{}, fmt.Errorf("error returned from OpenAI: %s", result.Error.Message)
	}

	return result, nil
}
