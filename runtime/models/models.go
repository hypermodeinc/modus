/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
	"github.com/hypermodeinc/modus/runtime/sentryutils"
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
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
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
			switch httpe.StatusCode {
			case http.StatusNotFound:
				if app.IsDevEnvironment() {
					return empty, fmt.Errorf("model %s is not available in the local dev environment", model.SourceModel)
				}
			case http.StatusUnauthorized:
				return empty, fmt.Errorf("invalid or expired API key. Please use `hyp login` to create a new API key")
			case http.StatusForbidden:
				return empty, fmt.Errorf("API key is disabled or usage limit exceeded")
			case http.StatusTooManyRequests:
				return empty, fmt.Errorf("rate limit exceeded")
			}
		}

		if res == nil {
			return empty, err
		}
	}

	// NOTE: This path occurs whether or not there's an error, as long as there was some response body content.
	db.WriteInferenceHistory(ctx, model, payload, res.Data, res.StartTime, res.EndTime)
	return res.Data, err
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
