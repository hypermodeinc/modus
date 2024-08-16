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

func (dc *dgraphConnector) alterSchema(ctx context.Context, schema string) (string, error) {
	op := &api.Operation{Schema: schema}
	if err := dc.dgClient.Alter(ctx, op); err != nil {
		return "", err
	}

	return "success", nil
}

func (dc *dgraphConnector) dropAttr(ctx context.Context, attr string) (string, error) {
	op := &api.Operation{DropAttr: attr}
	if err := dc.dgClient.Alter(ctx, op); err != nil {
		return "", err
	}

	return "success", nil
}

func (dc *dgraphConnector) dropAll(ctx context.Context) (string, error) {
	op := &api.Operation{DropAll: true}
	if err := dc.dgClient.Alter(ctx, op); err != nil {
		return "", err
	}

	return "success", nil
}

func (dc *dgraphConnector) executeQuery(ctx context.Context, query string, params map[string]string) (string, error) {
	if query == "" {
		return "", nil
	}

	tx := dc.dgClient.NewTxn()
	defer func() {
		if err := tx.Discard(ctx); err != nil {
			logger.Warn(ctx).Err(err).Msg("Error discarding transaction.")
			return
		}

	}()

	req := &api.Request{Query: query}

	if params != nil && len(params) > 0 {
		req.Vars = params
	}

	resp, err := tx.Do(ctx, req)
	if err != nil {
		return "", err
	}

	return string(resp.Json), nil
}

func (dc *dgraphConnector) executeMutations(ctx context.Context, mutations []string) (map[string]string, error) {
	if len(mutations) == 0 {
		return nil, nil
	}

	tx := dc.dgClient.NewTxn()
	defer func() {
		if err := tx.Discard(ctx); err != nil {
			logger.Warn(ctx).Err(err).Msg("Error discarding transaction.")
			return
		}

	}()

	mus := make([]*api.Mutation, 0, len(mutations))
	for _, m := range mutations {
		mus = append(mus, &api.Mutation{SetNquads: []byte(m)})
	}

	req := &api.Request{Mutations: mus, CommitNow: true}

	resp, err := tx.Do(ctx, req)
	return resp.Uids, err
}
