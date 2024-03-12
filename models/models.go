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
const HypermodeHost string = "hypermode"

func GetModel(modelName string, task config.ModelTask) (config.Model, error) {
	for _, model := range config.HypermodeData.Models {
		if model.Name == modelName && model.Task == task {
			return model, nil
		}
	}

	return config.Model{}, fmt.Errorf("a model '%s' for task '%s' was not found", modelName, task)
}

func GetModelKey(ctx context.Context, model config.Model) (string, error) {

	// First see if the model key is in an environment variable.
	// Try well-known environment variables first, then fall back to the model-specific environment variable
	key := getWellKnownEnvironmentVariable(model)
	if key == "" {
		key = os.Getenv(modelKeyPrefix + strings.ToUpper(model.Name))
	}

	// If the model key is not in an environment variable, try to get it from AWS Secrets Manager.
	// Use the model key as the secret name.
	if aws.Enabled() {
		return aws.GetSecretString(ctx, model.Name)
	}

	if key == "" {
		return "", fmt.Errorf("environment variable not found for key of model '%s'", model.Name)
	}
	return key, nil
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
	// self hosted models takes in array, can optimize for parallelizing later
	if model.Host == HypermodeHost {
		var sentences []string
		for _, v := range sentenceMap {
			sentences = append(sentences, v)
		}
		return utils.PostHttp[TResult](model.Endpoint, sentences, headers)
	}

	return utils.PostHttp[TResult](model.Endpoint, sentenceMap, headers)
}
