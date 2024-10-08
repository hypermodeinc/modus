/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package legacymodels

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/hosts"
	"github.com/hypermodeinc/modus/runtime/models"
	"github.com/hypermodeinc/modus/runtime/utils"
)

type outputFormat string

const (
	outputFormatText outputFormat = "text"
	outputFormatJson outputFormat = "json_object"
)

func InvokeTextGenerator(ctx context.Context, modelName string, format string, instruction string, sentence string) (string, error) {
	model, err := models.GetModel(modelName)
	if err != nil {
		return "", err
	}

	if model.Host != hosts.OpenAIHost {
		return "", fmt.Errorf("unsupported model host: %s", model.Host)
	}

	host, err := hosts.GetHttpHost(model.Host)
	if err != nil {
		return "", err
	}

	outputFormat := outputFormat(format)
	if outputFormat != outputFormatText && outputFormat != outputFormatJson {
		return "", fmt.Errorf("unsupported output format: %s", format)
	}

	result, err := callOpenAIChatCompletion(ctx, model, host, instruction, sentence, outputFormat)
	if err != nil {
		return "", err
	}

	if result.Error.Message != "" {
		err := fmt.Errorf("error received: %s", result.Error.Message)
		return "", err
	}

	if outputFormat == outputFormatJson {
		for _, choice := range result.Choices {
			_, err := utils.JsonSerialize(choice.Message.Content)
			if err != nil {
				return "", fmt.Errorf("one of the generated message is not a valid JSON: %w", err)
			}
		}
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty result returned from OpenAI")
	}

	return result.Choices[0].Message.Content, nil
}
