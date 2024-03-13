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

const modelKeyPrefix = "HYP_MODEL_KEY_"

func GetModel(modelName string, task config.ModelTask) (config.Model, error) {
	for _, model := range config.HypermodeData.Models {
		if model.Name == modelName && model.Task == task {
			return model, nil
		}
	}

	return config.Model{}, fmt.Errorf("a model '%s' for task '%s' was not found", modelName, task)
}

func GetModelKey(ctx context.Context, model config.Model) (string, error) {
	var key string
	var err error

	if config.UseAwsSecrets {
		// Get the model key from AWS Secrets Manager, using the model name as the secret.
		key, err = aws.GetSecretString(ctx, model.Name)
		if key != "" {
			return key, nil
		}
	} else {
		// Try well-known environment variables first, then model-specific environment variables.
		key := getWellKnownEnvironmentVariable(model)
		if key != "" {
			return key, nil
		}

		keyEnvVar := modelKeyPrefix + strings.ToUpper(model.Name)
		key = os.Getenv(keyEnvVar)
		if key != "" {
			return key, nil
		} else {
			err = fmt.Errorf("environment variable '%s' not found", keyEnvVar)
		}
	}

	return "", fmt.Errorf("error getting key for model '%s': %w", model.Name, err)
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
