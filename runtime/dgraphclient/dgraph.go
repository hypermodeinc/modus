/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraphclient

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/logger"

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

func (dc *dgraphConnector) execute(ctx context.Context, req *Request) (*Response, error) {
	if len(req.Mutations) == 0 {
		if req.Query.Query == "" {
			return nil, fmt.Errorf("no query or mutations provided")
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
			return nil, err
		}

		return &Response{Json: string(resp.Json), Uids: resp.Uids}, nil
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

	dgoReq := &api.Request{Mutations: mus, CommitNow: true}

	if req.Query.Query != "" {
		dgoReq.Query = req.Query.Query

		if len(req.Query.Variables) > 0 {
			dgoReq.Vars = req.Query.Variables
		}
	}

	resp, err := tx.Do(ctx, dgoReq)
	if err != nil {
		return nil, err
	}

	return &Response{Json: string(resp.Json), Uids: resp.Uids}, nil
}
