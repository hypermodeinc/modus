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

func getModelEndpointAndHost(model manifest.ModelInfo) (string, manifest.HostInfo, error) {

	host, err := hosts.GetHost(model.Host)
	if err != nil {
		return "", manifest.HostInfo{}, err
	}

	if host.Name == hosts.HypermodeHost {
		// TODO: make sure this is the correct endpoint URL for Kserve
		endpoint := fmt.Sprintf("http://%s.%s/", strings.ToLower(model.Name), config.ModelHost)
		return endpoint, host, nil
	}

	if host.BaseURL != "" && host.Endpoint != "" {
		return "", host, fmt.Errorf("specify either base URL or endpoint for a host, not both")
	}

	if host.BaseURL != "" {
		if model.Path == "" {
			return "", host, fmt.Errorf("model path is not defined")
		}
		endpoint := fmt.Sprintf("%s/%s", strings.TrimRight(host.BaseURL, "/"), strings.TrimLeft(model.Path, "/"))
		return endpoint, host, nil
	}

	if model.Path != "" {
		return "", host, fmt.Errorf("model path is defined but host has no base URL")
	}

	return host.Endpoint, host, nil
}

func PostToModelEndpoint[TResult any](ctx context.Context, model manifest.ModelInfo, payload any) (TResult, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	endpoint, host, err := getModelEndpointAndHost(model)
	if err != nil {
		var empty TResult
		return empty, err
	}

	var bs func(context.Context, *http.Request) error
	if host.Name != hosts.HypermodeHost {
		bs = func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Content-Type", "application/json")
			return secrets.ApplyHostSecrets(ctx, host, req)
		}
	}

	res, err := utils.PostHttp[TResult](ctx, endpoint, payload, bs)
	if err != nil {
		var empty TResult
		return empty, err
	}

	db.WriteInferenceHistory(ctx, model, payload, res.Data, res.StartTime, res.EndTime)

	return res.Data, nil
}

func PostToModelEndpoint_Old[TResult any](ctx context.Context, sentenceMap map[string]string, model manifest.ModelInfo) (map[string]TResult, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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

	res, err := PostToModelEndpoint[PredictionResult[TResult]](ctx, model, req)
	if err != nil {
		return nil, err
	}

	if len(res.Predictions) != len(keys) {
		return nil, fmt.Errorf("number of predictions does not match number of sentences")
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
