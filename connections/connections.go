package connections

import (
	"context"
	"fmt"
	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/utils"

	"github.com/buger/jsonparser"
	"github.com/hypermodeAI/manifest"
)

type request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func ExecuteGraphqlApi(ctx context.Context, host manifest.HostInfo, stmt string, vars map[string]any) (string, error) {
	request := request{
		Query:     stmt,
		Variables: vars,
	}

	headers := map[string]string{}

	if host.Endpoint == "" {
		return "", fmt.Errorf("host endpoint is not defined")
	}
	if host.AuthHeader != "" {
		key, err := hosts.GetHostKey(ctx, host)
		if err != nil {
			return "", err
		}

		headers[host.AuthHeader] = key
	}

	response, err := utils.PostHttp[[]byte](host.Endpoint, request, headers)
	if err != nil {
		return "", fmt.Errorf("error posting GraphQL statement: %w", err)
	}

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
