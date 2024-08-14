/*
 * Copyright 2024 Hypermode, Inc.
 */

package dqlclient

import (
	"context"
	"hmruntime/logger"

	"github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
	"google.golang.org/grpc"
)

type dgraphConnector struct {
	conn     *grpc.ClientConn
	dgClient *dgo.Dgraph
}

func (dc *dgraphConnector) execute(ctx context.Context, query string, mutations []string, params map[string]string) (string, error) {
	if query == "" && len(mutations) == 0 {
		return "", nil
	}

	tx := dc.dgClient.NewTxn()
	defer func() {
		if err := tx.Discard(ctx); err != nil {
			logger.Warn(ctx).Err(err).Msg("Error discarding transaction.")
			return
		}

	}()

	req := &api.Request{}

	mus := make([]*api.Mutation, 0, len(mutations))
	for _, m := range mutations {
		mus = append(mus, &api.Mutation{SetJson: []byte(m)})
	}

	if len(mus) > 0 {
		req.Mutations = mus
	}

	if query != "" {
		req.Query = query
	} else {
		req.CommitNow = true
	}

	if params != nil {
		req.Vars = params
	}

	resp, err := tx.Do(ctx, req)
	if err != nil {
		return "", err
	}

	return string(resp.Json), nil
}
