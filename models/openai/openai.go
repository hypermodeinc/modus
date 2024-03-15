/*
 * Copyright 2024 Hypermode, Inc.
 */

package openai

import (
	"context"
	"fmt"

	"hmruntime/config"
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

func GenerateText(ctx context.Context, model config.Model, instruction string, sentence string) (models.ChatResponse, error) {
	return ChatCompletion(ctx, model, instruction, sentence, "text")
}
func GenerateJson(ctx context.Context, model config.Model, instruction string, sentence string) (models.ChatResponse, error) {
	return ChatCompletion(ctx, model, instruction, sentence, "json_object")
}
func ChatCompletion(ctx context.Context, model config.Model, instruction string, sentence string, outputType string) (models.ChatResponse, error) {

	// Get the OpenAI API key to use for this model
	key, err := models.GetModelKey(ctx, model)
	if err != nil {
		return models.ChatResponse{}, err
	}

	// build the request body following OpenAI API
	reqBody := ChatContext{
		Model: model.SourceModel,
		ResponseFormat: ResponseFormat{
			Type: outputType,
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

	return result, nil
}
