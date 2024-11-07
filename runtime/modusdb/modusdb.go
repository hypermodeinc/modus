/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package modusdb

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modusdb"
)

var mdbExt *modusdb.DB
var mdbInt *modusdb.DB
var dataDir = "data"
var internalDataDir = "internal_data"

func Init(ctx context.Context) {
	var err error
	mdbExt, err = modusdb.New(modusdb.NewDefaultConfig().WithDataDir(dataDir))
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb")
	}
	mdbInt, err = modusdb.New(modusdb.NewDefaultConfig().WithDataDir(internalDataDir))
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb")
	}
	err = mdbInt.AlterSchema(ctx, internalSchema)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb schema")
	}
}

func Close(ctx context.Context) {
	if mdbExt != nil {
		mdbExt.Close()
	}
	if mdbInt != nil {
		mdbInt.Close()
	}
}

func ExtDropAll(ctx context.Context) (string, error) {
	err := mdbExt.DropAll(ctx)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func ExtDropData(ctx context.Context) (string, error) {
	err := mdbExt.DropData(ctx)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func ExtAlterSchema(ctx context.Context, schema string) (string, error) {
	totalSchema := fmt.Sprintf("%s\n%s", internalSchema, schema)
	err := mdbExt.AlterSchema(ctx, totalSchema)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func ExtMutate(ctx context.Context, mutationReq MutationRequest) (map[string]uint64, error) {
	mutations := mutationReq.Mutations
	mus := make([]*api.Mutation, 0, len(mutations))
	for _, m := range mutations {
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
	return mdbExt.Mutate(ctx, mus)
}

func ExtQuery(ctx context.Context, query string) (*Response, error) {
	resp, err := mdbExt.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	hdrs := make(map[string]*ListOfString, len(resp.Hdrs))
	for k, v := range resp.Hdrs {
		hdrs[k] = &ListOfString{Value: v.Value}
	}
	return &Response{
		Json: string(resp.Json),
		Txn: &TxnContext{
			StartTs:  resp.Txn.StartTs,
			CommitTs: resp.Txn.CommitTs,
			Aborted:  resp.Txn.Aborted,
			Keys:     resp.Txn.Keys,
			Preds:    resp.Txn.Preds,
			Hash:     resp.Txn.Hash,
		},
		Latency: &Latency{
			ParsingNs:         resp.Latency.ParsingNs,
			ProcessingNs:      resp.Latency.ProcessingNs,
			EncodingNs:        resp.Latency.EncodingNs,
			AssignTimestampNs: resp.Latency.AssignTimestampNs,
			TotalNs:           resp.Latency.TotalNs,
		},
		Metrics: &Metrics{
			NumUids: resp.Metrics.NumUids,
		},
		Uids: resp.Uids,
		Rdf:  resp.Rdf,
		Hdrs: hdrs,
	}, nil
}

func IntMutate(ctx context.Context, mutationReq MutationRequest) (map[string]uint64, error) {
	mutations := mutationReq.Mutations
	mus := make([]*api.Mutation, 0, len(mutations))
	for _, m := range mutations {
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
	return mdbInt.Mutate(ctx, mus)
}

func IntQuery(ctx context.Context, query string) (*Response, error) {
	resp, err := mdbInt.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	hdrs := make(map[string]*ListOfString, len(resp.Hdrs))
	for k, v := range resp.Hdrs {
		hdrs[k] = &ListOfString{Value: v.Value}
	}
	return &Response{
		Json: string(resp.Json),
		Txn: &TxnContext{
			StartTs:  resp.Txn.StartTs,
			CommitTs: resp.Txn.CommitTs,
			Aborted:  resp.Txn.Aborted,
			Keys:     resp.Txn.Keys,
			Preds:    resp.Txn.Preds,
			Hash:     resp.Txn.Hash,
		},
		Latency: &Latency{
			ParsingNs:         resp.Latency.ParsingNs,
			ProcessingNs:      resp.Latency.ProcessingNs,
			EncodingNs:        resp.Latency.EncodingNs,
			AssignTimestampNs: resp.Latency.AssignTimestampNs,
			TotalNs:           resp.Latency.TotalNs,
		},
		Metrics: &Metrics{
			NumUids: resp.Metrics.NumUids,
		},
		Uids: resp.Uids,
		Rdf:  resp.Rdf,
		Hdrs: hdrs,
	}, nil
}
