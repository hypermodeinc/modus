/*
 * Copyright 2024 Hypermode, Inc.
 */

package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/db"
	"hmruntime/hosts"
	"hmruntime/models"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
)

type ChatContext struct {
	Model          string               `json:"model"`
	ResponseFormat ResponseFormat       `json:"response_format"`
	Messages       []models.ChatMessage `json:"messages"`
}
type ResponseFormat struct {
	Type string `json:"type"`
}

func ChatCompletion(ctx context.Context, model manifest.ModelInfo, host manifest.HostInfo, instruction string, sentence string, outputFormat models.OutputFormat) (models.ChatResponse, error) {

	// Get the OpenAI API key to use for this model
	key, err := hosts.GetHostKey(ctx, host)
	if err != nil {
		return models.ChatResponse{}, err
	}

	// build the request body following OpenAI API
	reqBody := ChatContext{
		Model: model.SourceModel,
		ResponseFormat: ResponseFormat{
			Type: string(outputFormat),
		},
		Messages: []models.ChatMessage{
			{Role: "system", Content: instruction},
			{Role: "user", Content: sentence},
		},
	}

	// We ignore the model endpoint and use the OpenAI endpoint
	const endpoint = "https://api.openai.com/v1/chat/completions"
	headers := map[string]string{
		"Authorization": "Bearer " + key,
	}

	start := utils.GetTime()
	result, err := utils.PostHttp[models.ChatResponse](endpoint, reqBody, headers)
	end := utils.GetTime()

	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("error posting to OpenAI: %w", err)
	}

	if result.Error.Message != "" {
		return models.ChatResponse{}, fmt.Errorf("error returned from OpenAI: %s", result.Error.Message)
	}

	// write the results to the database
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return result, fmt.Errorf("error marshalling result: %w", err)
	}
	db.WriteInferenceHistory(ctx, model, sentence, string(resultBytes), start, end)

	return result, nil
}
