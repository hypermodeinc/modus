/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/utils"
)

func GetModelSpec(modelName string, modelType config.ModelType) (config.ModelSpec, error) {
	for _, modelSpec := range config.HypermodeData.ModelSpecs {
		if modelSpec.Name == modelName && modelSpec.ModelType == modelType {
			return modelSpec, nil
		}
	}

	return config.ModelSpec{}, fmt.Errorf("a model '%s' of type '%s' was not found", modelName, modelType)
}

func GetModelKey(ctx context.Context, modelSpec config.ModelSpec) (string, error) {
	if aws.UseAwsForPluginStorage() {
		return aws.GetSecretString(ctx, modelSpec.Name)
	} else {
		return modelSpec.ApiKey, nil
	}
}

func PostToModelEndpoint[TResult any](ctx context.Context, sentenceMap map[string]string, modelSpec config.ModelSpec) (TResult, error) {
	key, err := GetModelKey(ctx, modelSpec)
	if err != nil {
		var result TResult
		return result, fmt.Errorf("error getting model key secret: %w", err)
	}

	headers := map[string]string{
		modelSpec.AuthHeader: key,
	}

	return utils.PostHttp[TResult](modelSpec.Endpoint, sentenceMap, headers)
}
