/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"
	"os"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/utils"
)

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
		return aws.GetSecretString(ctx, model.Name)
	} else {
		// get key from env
		key := os.Getenv(model.Name)
		if key == "" {
			return "", fmt.Errorf("model key '%s' not found in environment", model.Name)
		}
		return key, nil
	}
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
