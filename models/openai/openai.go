/*
 * Copyright 2024 Hypermode, Inc.
 */

package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/config"
	"hmruntime/models"
	"hmruntime/utils"
)

type ChatContext struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
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

func GenerateText(ctx context.Context, modelSpec config.ModelSpec, instruction string, sentence string) (ChatResponse, error) {

	// Get the OpenAI API key to use for this model
	key, err := models.GetModelKey(ctx, modelSpec)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("error getting model key secret: %w", err)
	}

	// Encode the instruction and sentence as JSON strings
	jInstruction, _ := json.Marshal(instruction)
	jSentence, _ := json.Marshal(sentence)

	// build the request body following OpenAI API
	reqBody := ChatContext{
		Model: modelSpec.BaseModel,
		Messages: []ChatMessage{
			{Role: "system", Content: string(jInstruction)},
			{Role: "user", Content: string(jSentence)},
		},
	}

	// We ignore modelSpec.endpoint and use the OpenAI endpoint
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
