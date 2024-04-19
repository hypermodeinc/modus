package connections

import (
	"context"
	"fmt"
	"hmruntime/appdata"
	"hmruntime/hosts"
	"hmruntime/utils"
)

type request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func ExecuteGraphqlApi[TResponse any](ctx context.Context, host appdata.Host, stmt string, vars map[string]any) (TResponse, error) {
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

	response, err := utils.PostHttp[TResponse](host.Endpoint, request, nil)
	if err != nil {
		return response, fmt.Errorf("error posting GraphQL statement: %w", err)
	}

	return response, nil
}
