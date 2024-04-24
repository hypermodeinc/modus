package connections

import (
	"context"
	"fmt"
	"hmruntime/hosts"
	"hmruntime/manifest"
	"hmruntime/utils"
)

type request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func ExecuteGraphqlApi[TResponse any](ctx context.Context, host manifest.Host, stmt string, vars map[string]any) (TResponse, error) {
	request := request{
		Query:     stmt,
		Variables: vars,
	}

	var response TResponse
	headers := map[string]string{}

	if host.Endpoint == "" {
		return response, fmt.Errorf("host endpoint is not defined")
	}
	if host.AuthHeader != "" {
		key, err := hosts.GetHostKey(ctx, host)
		if err != nil {
			return response, err
		}

		headers[host.AuthHeader] = key
	}

	response, err := utils.PostHttp[TResponse](host.Endpoint, request, nil)
	if err != nil {
		return response, fmt.Errorf("error posting GraphQL statement: %w", err)
	}

	return response, nil
}

func ExecuteDQL[TResponse any](ctx context.Context, host manifest.Host, stmt string, vars map[string]any, isMutation bool) (TResponse, error) {
	request := request{
		Query:     stmt,
		Variables: vars,
	}

	var response TResponse
	headers := map[string]string{}

	if host.Endpoint == "" {
		return response, fmt.Errorf("host endpoint is not defined")
	}
	if host.AuthHeader != "" {
		key, err := hosts.GetHostKey(ctx, host)
		if err != nil {
			return response, fmt.Errorf("error getting model key secret: %w", err)
		}

		headers[host.AuthHeader] = key
	}

	var url string
	if isMutation {
		url = host.Endpoint + "/mutate?commitNow=true"
	} else {
		url = host.Endpoint + "/query"
	}

	response, err := utils.PostHttp[TResponse](url, request, nil)
	if err != nil {
		return response, fmt.Errorf("error posting DQL statement: %w", err)
	}

	return response, nil
}
