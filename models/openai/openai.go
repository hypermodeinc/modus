/*
 * Copyright 2024 Hypermode, Inc.
 */

package openai

import (
	"context"
	"fmt"
	"regexp"

	"hmruntime/db"
	"hmruntime/hosts"
	"hmruntime/manifestdata"
	"hmruntime/models"

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

var authHeaderRegex = regexp.MustCompile(`^Bearer {{\s*\w+?\s*}}$`)

func ChatCompletion(ctx context.Context, model manifest.ModelInfo, host manifest.HTTPHostInfo,
	instruction string, sentence string, outputFormat models.OutputFormat) (*models.ChatResponse, error) {

	// We ignore the model endpoint and use the OpenAI chat completion endpoint
	host.Endpoint = "https://api.openai.com/v1/chat/completions"

	// Validate that the host has an Authorization header containing a Bearer token
	authHeaderValue, ok := host.Headers["Authorization"]
	if !ok {
		// Handle how the key was passed in the old V1 manifest format
		if manifestdata.Manifest.Version == 1 && len(host.Headers) == 1 {
			for k, v := range host.Headers {
				apiKey := v
				delete(host.Headers, k)
				host.Headers["Authorization"] = "Bearer " + apiKey
			}
		} else {
			return nil, fmt.Errorf("host must have an Authorization header containing a Bearer token")
		}
	} else if !authHeaderRegex.MatchString(authHeaderValue) {
		return nil, fmt.Errorf("host Authorization header must be of the form: \"Bearer {{SECRET_NAME}}\"")
	}

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

	result, err := hosts.PostToHostEndpoint[models.ChatResponse](ctx, host, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error posting to OpenAI: %w", err)
	}

	if result.Data.Error.Message != "" {
		return &result.Data, fmt.Errorf("error returned from OpenAI: %s", result.Data.Error.Message)
	}

	db.WriteInferenceHistory(ctx, model, reqBody, result.Data, result.StartTime, result.EndTime)

	return &result.Data, nil
}
