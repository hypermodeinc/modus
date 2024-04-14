/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"
	"os"
	"strings"

	"hmruntime/appdata"
	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/utils"
)

const modelKeyPrefix = "HYP_MODEL_KEY_"
const HypermodeHost string = "hypermode"

// generic output format for models functions
// can be extended to support more formats
// for now, we support text and json_object used in generateText function
type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJson OutputFormat = "json_object"
)

func GetModel(modelName string, task appdata.ModelTask) (appdata.Model, error) {
	for _, model := range appdata.HypermodeData.Models {
		if model.Name == modelName {
			if (model.Task == task) || (model.Host == OpenAIHost) || (model.Host == MistralHost) {
				return model, nil
			}
		}
	}

	return appdata.Model{}, fmt.Errorf("a model '%s' for task '%s' was not found", modelName, task)
}

func GetModelKey(ctx context.Context, model appdata.Model) (string, error) {
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

func getWellKnownEnvironmentVariable(model appdata.Model) string {

	// Some model hosts have well-known environment variables that are used to store the model key.
	// We should support these to make it easier for users to set up their environment.
	// We can expand this list as we add more model hosts.

	switch model.Host {
	case OpenAIHost:
		return os.Getenv("OPENAI_API_KEY")
	}
	return ""
}

type PredictionResult[T any] struct {
	Predictions []T `json:"predictions"`
}

func PostToModelEndpoint[TResult any](ctx context.Context, sentenceMap map[string]string, model appdata.Model) (map[string]TResult, error) {
	// self hosted models takes in array, can optimize for parallelizing later
	keys, sentences := []string{}, []string{}

	for k, v := range sentenceMap {
		// create a map of sentences to send to the model
		sentences = append(sentences, v)
		// create a list of keys to map the results back to the original sentences
		keys = append(keys, k)
	}
	// create a map of sentences to send to the model
	req := map[string][]string{"instances": sentences}

	var endpoint string
	headers := map[string]string{}

	switch model.Host {
	case HypermodeHost:
		endpoint = fmt.Sprintf("http://%s.%s/%s:predict", model.Name, config.ModelHost, model.Task)
	default:
		// If the model is not hosted by Hypermode, we need to get the model key and add it to the request headers
		endpoint = model.Endpoint
		key, err := GetModelKey(ctx, model)
		if err != nil {
			return map[string]TResult{}, fmt.Errorf("error getting model key secret: %w", err)
		}

		headers[model.AuthHeader] = key
	}

	res, err := utils.PostHttp[PredictionResult[TResult]](endpoint, req, headers)
	if err != nil {
		return map[string]TResult{}, err
	}
	if len(res.Predictions) != len(keys) {
		return map[string]TResult{}, fmt.Errorf("number of predictions does not match number of sentences")
	}

	// map the results back to the original sentences
	result := make(map[string]TResult)
	for i, v := range res.Predictions {
		result[keys[i]] = v
	}
	return result, nil
}

// Define  structures used by text generation functions
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
