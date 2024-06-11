/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"hmruntime/config"
	"hmruntime/db"
	"hmruntime/hosts"
	"hmruntime/manifestdata"
	"hmruntime/secrets"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
)

// generic output format for models functions
// can be extended to support more formats
// for now, we support text and json_object used in generateText function
type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJson OutputFormat = "json_object"
)

func GetModel(modelName string) (manifest.ModelInfo, error) {
	model, ok := manifestdata.Manifest.Models[modelName]
	if !ok {
		return manifest.ModelInfo{}, fmt.Errorf("model %s was not found", modelName)
	}

	return model, nil
}

func GetModelForTask(modelName string, task manifest.ModelTask) (manifest.ModelInfo, error) {

	model, err := GetModel(modelName)
	if err != nil {
		return manifest.ModelInfo{}, err
	}

	if model.Task != task {
		return manifest.ModelInfo{}, fmt.Errorf("model %s is not a %s model", modelName, task)
	}

	return model, nil
}

func PostToModelEndpoint[TResult any](ctx context.Context, sentenceMap map[string]string, model manifest.ModelInfo) (map[string]TResult, error) {

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
	var bs func(context.Context, *http.Request) error

	switch model.Host {
	case hosts.HypermodeHost:
		endpoint = fmt.Sprintf("http://%s.%s/%s:predict", strings.ToLower(model.Name), config.ModelHost, model.Task)
	default:

		host, err := hosts.GetHost(model.Host)
		if err != nil {
			return nil, err
		}

		if host.Endpoint == "" {
			return nil, fmt.Errorf("host endpoint is not defined")
		}

		endpoint = host.Endpoint

		bs = func(ctx context.Context, req *http.Request) error {
			return secrets.ApplyHostSecrets(ctx, host, req)
		}
	}

	res, err := utils.PostHttp[PredictionResult[TResult]](ctx, endpoint, req, bs)
	if err != nil {
		return map[string]TResult{}, err
	}
	if len(res.Data.Predictions) != len(keys) {
		return map[string]TResult{}, fmt.Errorf("number of predictions does not match number of sentences")
	}

	// map the results back to the original sentences
	result := make(map[string]TResult)
	for i, v := range res.Data.Predictions {
		result[keys[i]] = v
	}

	db.WriteInferenceHistory(ctx, model, sentenceMap, result, res.StartTime, res.EndTime)

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
