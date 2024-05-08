/*
 * Copyright 2024 Hypermode, Inc.
 */

package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/hosts"
	"hmruntime/manifest"
	"hmruntime/models"
	"hmruntime/utils"
)

type ChatContext struct {
	Model          string               `json:"model"`
	ResponseFormat ResponseFormat       `json:"response_format"`
	Messages       []models.ChatMessage `json:"messages"`
}
type ResponseFormat struct {
	Type string `json:"type"`
}

func ChatCompletion(ctx context.Context, model manifest.Model, host manifest.Host, instruction string, sentence string, outputFormat models.OutputFormat) (models.ChatResponse, error) {

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

	result, err := utils.PostHttp[models.ChatResponse](endpoint, reqBody, headers)

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
	models.WriteInferenceHistoryToDB(model, sentence, string(resultBytes))

	return result, nil
}
