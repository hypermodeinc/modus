/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"

	"hmruntime/config"
	"hmruntime/hosts"
	"hmruntime/manifest"
	"hmruntime/utils"
)

const HypermodeHost string = "hypermode"

// generic output format for models functions
// can be extended to support more formats
// for now, we support text and json_object used in generateText function
type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJson OutputFormat = "json_object"
)

func GetModel(modelName string, task manifest.ModelTask) (manifest.Model, error) {
	for _, model := range manifest.HypermodeData.Models {
		if model.Name == modelName && model.Task == task {
			return model, nil
		}
	}

	return manifest.Model{}, fmt.Errorf("a model '%s' for task '%s' was not found", modelName, task)
}

func PostToModelEndpoint[TResult any](ctx context.Context, sentenceMap map[string]string, model manifest.Model) (map[string]TResult, error) {

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
		host, err := hosts.GetHost(model.Host)
		if err != nil {
			return map[string]TResult{}, err
		}

		endpoint = host.Endpoint
		if host.AuthHeader == "" {
			break
		}
		key, err := hosts.GetHostKey(ctx, host)
		if err != nil {
			return map[string]TResult{}, err
		}

		headers[host.AuthHeader] = key
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

type PredictionResult[T any] struct {
	Predictions []T `json:"predictions"`
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
