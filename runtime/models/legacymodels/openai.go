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
	"regexp"

	"github.com/hypermodeinc/modus/pkg/manifest"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/hosts"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
)

var authHeaderRegex = regexp.MustCompile(`^Bearer {{\s*\w+?\s*}}$`)

func callOpenAIChatCompletion(ctx context.Context, model *manifest.ModelInfo, host *manifest.HTTPHostInfo, instruction string, sentence string, outputFormat outputFormat) (*chatResponse, error) {

	// We ignore the model endpoint and use the OpenAI chat completion endpoint
	host.Endpoint = "https://api.openai.com/v1/chat/completions"

	// Validate that the host has an Authorization header containing a Bearer token
	authHeaderValue, ok := host.Headers["Authorization"]
	if !ok {
		// Handle how the key was passed in the old V1 manifest format
		if manifestdata.GetManifest().Version == 1 && len(host.Headers) == 1 {
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

	reqBody := chatContext{
		Model: model.SourceModel,
		ResponseFormat: responseFormat{
			Type: string(outputFormat),
		},
		Messages: []chatMessage{
			{Role: "system", Content: instruction},
			{Role: "user", Content: sentence},
		},
	}

	result, err := hosts.PostToHostEndpoint[chatResponse](ctx, host, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error posting to OpenAI: %w", err)
	}

	if result.Data.Error.Message != "" {
		return &result.Data, fmt.Errorf("error returned from OpenAI: %s", result.Data.Error.Message)
	}

	db.WriteInferenceHistory(ctx, model, reqBody, result.Data, result.StartTime, result.EndTime)

	return &result.Data, nil
}

type chatContext struct {
	Model          string         `json:"model"`
	ResponseFormat responseFormat `json:"response_format"`
	Messages       []chatMessage  `json:"messages"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []messageChoice `json:"choices"`
	Error   invokeError     `json:"error"`
}

type messageChoice struct {
	Message chatMessage `json:"message"`
}

type invokeError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}
