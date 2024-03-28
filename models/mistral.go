/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"
	"hmruntime/appdata"
	"hmruntime/utils"
)

type mistral struct{}

func (llm *mistral) ChatCompletion(ctx context.Context, model appdata.Model, instruction string, sentence string, outputFormat OutputFormat) (ChatResponse, error) {

	// Get the API key to use for this model
	key, err := GetModelKey(ctx, model)
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
