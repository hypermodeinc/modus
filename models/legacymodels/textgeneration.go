/*
 * Copyright 2024 Hypermode, Inc.
 */

package legacymodels

import (
	"context"
	"fmt"

	"hmruntime/hosts"
	"hmruntime/models"
	"hmruntime/utils"
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
