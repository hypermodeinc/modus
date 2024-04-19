/*
 * Copyright 2023 Hypermode, Inc.
 */

package dgraph

import (
	"context"
	"fmt"

	"hmruntime/config"
	"hmruntime/utils"
)

type dgraphRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func ExecuteDQL[TResponse any](ctx context.Context, stmt string, vars map[string]any, isMutation bool) (TResponse, error) {
	var url string
	if isMutation {
		url = config.DgraphUrl + "/mutate?commitNow=true"
	} else {
		url = config.DgraphUrl + "/query"
	}

	request := dgraphRequest{
		Query:     stmt,
		Variables: vars,
	}

	response, err := utils.PostHttp[TResponse](url, request, nil)
	if err != nil {
		return response, fmt.Errorf("error posting DQL statement: %w", err)
	}

	return response, nil
}

func ExecuteGQL[TResponse any](ctx context.Context, stmt string, vars map[string]any) (TResponse, error) {
	url := config.DgraphUrl + "/graphql"
	request := dgraphRequest{
		Query:     stmt,
		Variables: vars,
	}

	response, err := utils.PostHttp[TResponse](url, request, nil)
	if err != nil {
		return response, fmt.Errorf("error posting GraphQL statement: %w", err)
	}

	return response, nil
}
