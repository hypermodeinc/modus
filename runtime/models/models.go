/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package models

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/hosts"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/hypermodeinc/modus/runtime/utils"
)

var localHypermodeModels = map[string]bool{"meta-llama/meta-llama-3.1-8b-instruct": true, "sentence-transformers/all-minilm-l6-v2": true, "antoinemc/distilbart-mnli-github-issues": true, "distilbert/distilbert-base-uncased-finetuned-sst-2-english": true}

func GetModel(modelName string) (*manifest.ModelInfo, error) {
	model, ok := manifestdata.GetManifest().Models[modelName]
	if !ok {
		return nil, fmt.Errorf("model %s was not found", modelName)
	}

	return &model, nil
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

func PostToModelEndpoint[TResult any](ctx context.Context, model *manifest.ModelInfo, payload any) (TResult, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	host, err := hosts.GetHttpHost(model.Host)
	if err != nil {
		var empty TResult
		return empty, err
	}

	url, err := getModelEndpointUrl(model, host)
	if err != nil {
		var empty TResult
		return empty, err
	}

	bs := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Content-Type", "application/json")
		if host.Name != hosts.HypermodeHost {
			return secrets.ApplyHostSecretsToHttpRequest(ctx, host, req)
		}
		if config.IsDevEnvironment() {
			return secrets.ApplyAuthToLocalModelRequest(ctx, host, req)
		}
		return nil
	}

	res, err := utils.PostHttp[TResult](ctx, url, payload, bs)
	if err != nil {
		var empty TResult
		return empty, err
	}

	db.WriteInferenceHistory(ctx, model, payload, res.Data, res.StartTime, res.EndTime)

	return res.Data, nil
}

func getModelEndpointUrl(model *manifest.ModelInfo, host *manifest.HTTPHostInfo) (string, error) {

	if host.Name == hosts.HypermodeHost {
		if config.IsDevEnvironment() {
			if !isValidLocalHypermodeModel(model.SourceModel) {
				return "", fmt.Errorf("model %s is not available in the local Hypermode environment", model.SourceModel)
			}
			endpoint := fmt.Sprintf("https://models.hypermode.host/%s", strings.ToLower(model.SourceModel))
			return endpoint, nil
		}
		// Compose the Hypermode hosted model endpoint URL.
		// Example: http://modelname.bckid.svc.cluster.local/v1/models/modelname:predict
		endpoint := fmt.Sprintf("http://%s.%s/%[1]s:predict", strings.ToLower(model.Name), config.ModelHost)
		return endpoint, nil
	}

	if host.BaseURL != "" && host.Endpoint != "" {
		return "", fmt.Errorf("specify either base URL or endpoint for a host, not both")
	}

	if host.BaseURL != "" {
		if model.Path == "" {
			return "", fmt.Errorf("model path is not defined")
		}
		endpoint := fmt.Sprintf("%s/%s", strings.TrimRight(host.BaseURL, "/"), strings.TrimLeft(model.Path, "/"))
		return endpoint, nil
	}

	if model.Path != "" {
		return "", fmt.Errorf("model path is defined but host has no base URL")
	}

	return host.Endpoint, nil
}

func isValidLocalHypermodeModel(modelName string) bool {
	_, ok := localHypermodeModels[strings.ToLower(modelName)]
	return ok
}
