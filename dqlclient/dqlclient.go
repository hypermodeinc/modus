/*
 * Copyright 2024 Hypermode, Inc.
 */

package dqlclient

import (
	"context"
	"fmt"
	"hmruntime/utils"
)

func ExecuteQuery(ctx context.Context, hostName, query string, paramsJson string) (string, error) {
	var params map[string]string
	if err := utils.JsonDeserialize([]byte(paramsJson), &params); err != nil {
		return "", fmt.Errorf("error deserializing database query parameters: %w", err)
	}

	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return "", err
	}

	return dc.executeQuery(ctx, query, params)
}

func ExecuteMutations(ctx context.Context, hostName string, setMutations, delMutations []string) (map[string]string, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return nil, err
	}

	return dc.executeMutations(ctx, setMutations, delMutations)
}

func AlterSchema(ctx context.Context, hostName, schema string) (string, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return "", err
	}

	return dc.alterSchema(ctx, schema)
}

func DropAttr(ctx context.Context, hostName, attr string) (string, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return "", err
	}

	return dc.dropAttr(ctx, attr)
}

func DropAll(ctx context.Context, hostName string) (string, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return "", err
	}

	return dc.dropAll(ctx)
}
