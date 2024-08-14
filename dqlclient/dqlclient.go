/*
 * Copyright 2024 Hypermode, Inc.
 */

package dqlclient

import (
	"context"
	"fmt"
	"hmruntime/utils"
)

func ExecuteQuery(ctx context.Context, hostName, query string, mutations []string, paramsJson string) (string, error) {
	var params map[string]string
	if err := utils.JsonDeserialize([]byte(paramsJson), &params); err != nil {
		return "", fmt.Errorf("error deserializing database query parameters: %w", err)
	}

	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return "", err
	}

	return dc.execute(ctx, query, mutations, params)
}
