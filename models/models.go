/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"
	"os"
	"strings"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/utils"
)

const modelKeyPrefix = "HYPERMODE_MODEL_KEY_"

func GetModel(modelName string, task config.ModelTask) (config.Model, error) {
	for _, model := range config.HypermodeData.Models {
		if model.Name == modelName && model.Task == task {
			return model, nil
		}
	}

	return config.Model{}, fmt.Errorf("a model '%s' for task '%s' was not found", modelName, task)
}

func GetModelKey(ctx context.Context, model config.Model) (string, error) {
	if aws.Enabled() {
		// Get the model key from AWS Secrets Manager, using the model name as the secret name
		return aws.GetSecretString(ctx, model.Name)
	} else {
		// Try well-known environment variables first, then fall back to the model-specific environment variable
		key := getWellKnownEnvironmentVariable(model)
		if key == "" {
			key = os.Getenv(modelKeyPrefix + strings.ToUpper(model.Name))
		}
		if key == "" {
			return "", fmt.Errorf("environment variable not found for key of model '%s'", model.Name)
		}
		return key, nil
	}
}

func getWellKnownEnvironmentVariable(model config.Model) string {

	// Some model hosts have well-known environment variables that are used to store the model key.
	// We should support these to make it easier for users to set up their environment.
	// We can expand this list as we add more model hosts.

	switch model.Host {
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	}
	return ""
}

func PostToModelEndpoint[TResult any](ctx context.Context, sentenceMap map[string]string, model config.Model) (TResult, error) {
	key, err := GetModelKey(ctx, model)
	if err != nil {
		var result TResult
		return result, fmt.Errorf("error getting model key secret: %w", err)
	}

	headers := map[string]string{
		model.AuthHeader: key,
	}

	return utils.PostHttp[TResult](model.Endpoint, sentenceMap, headers)
}
