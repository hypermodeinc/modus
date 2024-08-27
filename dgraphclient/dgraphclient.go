/*
 * Copyright 2024 Hypermode, Inc.
 */

package dgraphclient

import (
	"context"
)

func Execute(ctx context.Context, hostName string, req Request) (Response, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return Response{}, err
	}

	return dc.execute(ctx, req)
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
