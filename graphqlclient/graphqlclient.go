/*
 * Copyright 2024 Hypermode, Inc.
 */

package graphqlclient

import (
	"context"
	"fmt"

	"hypruntime/hosts"
	"hypruntime/logger"
	"hypruntime/utils"

	"github.com/buger/jsonparser"
)

type graphqlRequestPayload struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func Execute(ctx context.Context, hostName string, stmt string, varsJson string) (string, error) {

	host, err := hosts.GetHttpHost(hostName)
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

	result, err := hosts.PostToHostEndpoint[[]byte](ctx, host, payload)
	if err != nil {
		return "", fmt.Errorf("error posting GraphQL statement: %w", err)
	}

	// Check for errors in the response so we can log them.
	response := result.Data
	gqlErrors, dataType, _, err := jsonparser.Get(response, "errors")
	if err != nil && err != jsonparser.KeyPathNotFoundError {
		return "", fmt.Errorf("error parsing GraphQL response: %w", err)
	}
	if dataType == jsonparser.Array && len(gqlErrors) > 0 {
		logger.Warn(ctx).
			Bool("user_visible", true).
			Str("errors", string(gqlErrors)).
			Msg("GraphQL API call returned errors.")
	}

	return string(response), nil
}
