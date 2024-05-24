/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/functions/assemblyscript"
	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/utils"

	"github.com/buger/jsonparser"
	"github.com/hypermodeAI/manifest"
	wasm "github.com/tetratelabs/wazero/api"
)

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pHostName uint32, pStmt uint32, pVars uint32) uint32 {

	var hostName, stmt, sVars string
	err := readParams3(ctx, mod, pHostName, pStmt, pVars, &hostName, &stmt, &sVars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	vars := make(map[string]any)
	if err := utils.JsonDeserialize([]byte(sVars), &vars); err != nil {
		logger.Err(ctx, err).Msg("Error deserializing GraphQL variables.")
		return 0
	}

	host, err := hosts.GetHost(hostName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting host.")
		return 0
	}

	result, err := executeGraphql(ctx, host, stmt, vars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing GraphQL operation.")
		return 0
	}

	offset, err := assemblyscript.WriteString(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func executeGraphql(ctx context.Context, host manifest.HostInfo, stmt string, vars map[string]any) (string, error) {
	request := graphqlRequest{
		Query:     stmt,
		Variables: vars,
	}

	headers := map[string]string{}

	if host.Endpoint == "" {
		return "", fmt.Errorf("host endpoint is not defined")
	}
	if host.AuthHeader != "" {
		key, err := hosts.GetHostSecret(ctx, host, host.AuthHeader)
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
