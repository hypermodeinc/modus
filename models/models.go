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

func GetModel(modelName string) (manifest.ModelInfo, error) {
	model, ok := manifestdata.GetManifest().Models[modelName]
	if !ok {
		return manifest.ModelInfo{}, fmt.Errorf("model %s was not found", modelName)
	}

	return model, nil
}

func GetModelInfo(modelName string) (*ModelInfo, error) {
	model, err := GetModel(modelName)
	if err != nil {
		return nil, err
	}

	info := &ModelInfo{
		Name:     model.Name,
		FullName: model.SourceModel,
	}

	return info, nil
}

func InvokeModel(ctx context.Context, modelName string, input string) (string, error) {
	model, err := GetModel(modelName)
	if err != nil {
		return "", err
	}

	// TODO: use the provider pattern instead of branching
	if model.Host == "aws-bedrock" {
		return invokeAwsBedrockModel(ctx, model, input)
	}

	return PostToModelEndpoint[string](ctx, model, input)
}

func PostToModelEndpoint[TResult any](ctx context.Context, model manifest.ModelInfo, payload any) (TResult, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	endpoint, host, err := getModelEndpointAndHost(model)
	if err != nil {
		var empty TResult
		return empty, err
	}

	bs := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Content-Type", "application/json")
		if host.Name != hosts.HypermodeHost {
			return secrets.ApplyHostSecretsToHttpRequest(ctx, host, req)
		}

		return nil
	}

	res, err := utils.PostHttp[TResult](ctx, endpoint, payload, bs)
	if err != nil {
		var empty TResult
		return empty, err
	}

	db.WriteInferenceHistory(ctx, model, payload, res.Data, res.StartTime, res.EndTime)

	return res.Data, nil
}

func getModelEndpointAndHost(model manifest.ModelInfo) (string, manifest.HTTPHostInfo, error) {

	host, err := hosts.GetHttpHost(model.Host)
	if err != nil {
		return "", manifest.HTTPHostInfo{}, err
	}

	if host.Name == hosts.HypermodeHost {
		// Compose the Hypermode hosted model endpoint URL.
		// Example: http://modelname.bckid.svc.cluster.local/v1/models/modelname:predict
		endpoint := fmt.Sprintf("http://%s.%s/%[1]s:predict", strings.ToLower(model.Name), config.ModelHost)
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
