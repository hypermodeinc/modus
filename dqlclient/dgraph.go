/*
 * Copyright 2024 Hypermode, Inc.
 */

package dqlclient

import (
	"context"

	"github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
	"google.golang.org/grpc"
)

type dgraphConnector struct {
	conn     *grpc.ClientConn
	dgClient *dgo.Dgraph
}

func (dc *dgraphConnector) query(ctx context.Context, stmt string, params map[string]string) (string, error) {
	tx := dc.dgClient.NewTxn()
	defer tx.Discard(ctx)

	req := &api.Request{
		Query: stmt,
		Vars:  params,
	}

	resp, err := tx.Do(ctx, req)
	if err != nil {
		return "", err
	}

	return string(resp.Json), nil
}
