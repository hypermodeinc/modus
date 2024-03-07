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
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func ExecuteDQL[TResponse any](ctx context.Context, stmt string, vars map[string]string, isMutation bool) (TResponse, error) {
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

func ExecuteGQL[TResponse any](ctx context.Context, stmt string, vars map[string]string) (TResponse, error) {
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

func GetGQLSchema(ctx context.Context) (string, error) {

	type DqlResponse[T any] struct {
		Data T `json:"data"`
	}

	type SchemaResponse struct {
		Node []struct {
			Schema string `json:"dgraph.graphql.schema"`
		} `json:"node"`
	}

	const query = "{node(func:has(dgraph.graphql.schema)){dgraph.graphql.schema}}"

	response, err := ExecuteDQL[DqlResponse[SchemaResponse]](ctx, query, nil, false)
	if err != nil {
		return "", fmt.Errorf("error getting GraphQL schema from Dgraph: %w", err)
	}

	data := response.Data
	if len(data.Node) == 0 {
		return "", fmt.Errorf("no GraphQL schema found in Dgraph")
	}

	return data.Node[0].Schema, nil
}
