/*
 * Copyright 2024 Hypermode, Inc.
 */

package dgraphclient

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

func (dc *dgraphConnector) execute(ctx context.Context, req Request) (Response, error) {
	if len(req.Mutations) == 0 {
		if req.Query.Query == "" {
			return Response{}, nil
		}

		tx := dc.dgClient.NewReadOnlyTxn()
		defer func() {
			if err := tx.Discard(ctx); err != nil {
				logger.Warn(ctx).Err(err).Msg("Error discarding transaction.")
				return
			}

		}()

		dgoReq := &api.Request{Query: req.Query.Query}

		if len(req.Query.Variables) > 0 {
			dgoReq.Vars = req.Query.Variables
		}

		resp, err := tx.Do(ctx, dgoReq)
		if err != nil {
			return Response{}, err
		}

		return Response{Json: string(resp.Json), Uids: resp.Uids}, nil
	}

	tx := dc.dgClient.NewTxn()
	defer func() {
		if err := tx.Discard(ctx); err != nil {
			logger.Warn(ctx).Err(err).Msg("Error discarding transaction.")
			return
		}

	}()

	mus := make([]*api.Mutation, 0, len(req.Mutations))

	for _, m := range req.Mutations {
		mu := &api.Mutation{}
		if m.SetJson != "" {
			mu.SetJson = []byte(m.SetJson)
		}
		if m.DelJson != "" {
			mu.DeleteJson = []byte(m.DelJson)
		}
		if m.SetNquads != "" {
			mu.SetNquads = []byte(m.SetNquads)
		}
		if m.DelNquads != "" {
			mu.DelNquads = []byte(m.DelNquads)
		}
		if m.Condition != "" {
			mu.Cond = m.Condition
		}

		mus = append(mus, mu)
	}

	dgoReq := &api.Request{Mutations: mus, CommitNow: req.CommitNow}

	if req.Query.Query != "" {
		dgoReq.Query = req.Query.Query

		if len(req.Query.Variables) > 0 {
			dgoReq.Vars = req.Query.Variables
		}
	}

	resp, err := tx.Do(ctx, dgoReq)
	if err != nil {
		return Response{}, err
	}

	return Response{Json: string(resp.Json), Uids: resp.Uids}, nil
}

func (dc *dgraphConnector) executeUpserts(ctx context.Context, query string, setMutations, delMutations []string) (map[string]string, error) {
	if len(setMutations) == 0 && len(delMutations) == 0 {
		return nil, nil
	}

	tx := dc.dgClient.NewTxn()
	defer func() {
		if err := tx.Discard(ctx); err != nil {
			logger.Warn(ctx).Err(err).Msg("Error discarding transaction.")
			return
		}

	}()

	mus := make([]*api.Mutation, 0, len(setMutations)+len(delMutations))

	for _, m := range setMutations {
		mus = append(mus, &api.Mutation{SetNquads: []byte(m)})
	}

	for _, m := range delMutations {
		mus = append(mus, &api.Mutation{DelNquads: []byte(m)})
	}

	req := &api.Request{Query: query, Mutations: mus, CommitNow: true}

	resp, err := tx.Do(ctx, req)
	return resp.Uids, err
}
