/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphqlclient

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/httpclient"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/tidwall/gjson"
)

type graphqlRequestPayload struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func ExecuteQuery(ctx context.Context, hostName string, stmt string, varsJson string) (string, error) {

	host, err := httpclient.GetHttpConnectionInfo(hostName)
	if err != nil {
		return "", err
	}

	vars := make(map[string]any)
	if err := utils.JsonDeserialize([]byte(varsJson), &vars); err != nil {
		return "", err
	}

	// https://graphql.org/learn/serving-over-http/
	payload := graphqlRequestPayload{
		Query:     stmt,
		Variables: vars,
	}

	result, err := httpclient.PostToConnectionEndpoint[[]byte](ctx, host, payload)
	if err != nil {
		return "", fmt.Errorf("error posting GraphQL statement: %w", err)
	}

	// Check if the response is valid JSON.
	if !gjson.ValidBytes(result.Data) {
		return "", fmt.Errorf("response from GraphQL API is not valid JSON: %s", string(result.Data))
	}

	// Check for errors in the response so we can log them.
	errorRes := gjson.GetBytes(result.Data, "errors")
	if errorRes.Exists() && errorRes.IsArray() && len(errorRes.Array()) > 0 {
		logger.Warn(ctx).
			Bool("user_visible", true).
			Str("errors", errorRes.String()).
			Msg("GraphQL API call returned errors.")
	}

	return string(result.Data), nil
}
