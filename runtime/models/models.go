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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/httpclient"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/hypermodeinc/modus/runtime/utils"
)

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

	// NOTE: Bedrock support is temporarily disabled
	// TODO: use the provider pattern instead of branching
	// if model.Connection == "aws-bedrock" {
	// 	return invokeAwsBedrockModel(ctx, model, input)
	// }

	return PostToModelEndpoint[string](ctx, model, input)
}

func PostToModelEndpoint[TResult any](ctx context.Context, model *manifest.ModelInfo, payload any) (TResult, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	connInfo, err := httpclient.GetHttpConnectionInfo(model.Connection)
	if err != nil {
		var empty TResult
		return empty, err
	}

	url, err := getModelEndpointUrl(model, connInfo)
	if err != nil {
		var empty TResult
		return empty, err
	}

	bs := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Content-Type", "application/json")
		if connInfo.Name == httpclient.HypermodeConnectionName {
			return authenticateHypermodeModelRequest(ctx, req, connInfo)
		} else {
			return secrets.ApplySecretsToHttpRequest(ctx, connInfo, req)
		}
	}

	res, err := utils.PostHttp[TResult](ctx, url, payload, bs)
	if err != nil {
		var empty TResult
		var httpe *utils.HttpError
		if errors.As(err, &httpe) {
			if app.IsDevEnvironment() && httpe.StatusCode == http.StatusNotFound {
				return empty, fmt.Errorf("model %s is not available in the local dev environment", model.SourceModel)
			}
		}

		return empty, err
	}

	db.WriteInferenceHistory(ctx, model, payload, res.Data, res.StartTime, res.EndTime)

	return res.Data, nil
}

func getModelEndpointUrl(model *manifest.ModelInfo, connection *manifest.HTTPConnectionInfo) (string, error) {

	if connection.Name == httpclient.HypermodeConnectionName {
		return getHypermodeModelEndpointUrl(model)
	}

	if connection.BaseURL != "" && connection.Endpoint != "" {
		return "", fmt.Errorf("specify either base URL or endpoint for an HTTP connection, not both")
	}

	if connection.BaseURL != "" {
		if model.Path == "" {
			return "", fmt.Errorf("model path is not defined")
		}
		endpoint := fmt.Sprintf("%s/%s", strings.TrimRight(connection.BaseURL, "/"), strings.TrimLeft(model.Path, "/"))
		return endpoint, nil
	}

	if model.Path != "" {
		return "", fmt.Errorf("model path is defined but the connection has no base URL")
	}

	return connection.Endpoint, nil
}
