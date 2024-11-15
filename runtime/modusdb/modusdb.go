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

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modusdb"
)

var mdbConfig modusdb.Config
var db *modusdb.DB
var mdbNsExt *modusdb.Namespace
var mdbNsInt *modusdb.Namespace
var dataDir = "data"

func Init(ctx context.Context) {
	var err error
	mdbConfig = modusdb.NewDefaultConfig(dataDir)
	db, err = modusdb.New(mdbConfig)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb")
	}

	mdbNsExt, err = db.GetNamespace(0)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb")
	}

	mdbNsInt, err = db.GetNamespace(1)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb namespace")
	}
	err = mdbNsInt.AlterSchema(ctx, internalSchema)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb schema")
	}
}

func Close(ctx context.Context) {
	if db != nil {
		db.Close()
	}
}

func DropData(ctx context.Context) (string, error) {
	err := mdbNsExt.DropData(ctx)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func AlterSchema(ctx context.Context, schema string) (string, error) {
	err := mdbNsExt.AlterSchema(ctx, schema)
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
	return mdbNsExt.Mutate(ctx, mus)
}

func ExtQuery(ctx context.Context, query string) (*Response, error) {
	resp, err := mdbNsExt.Query(ctx, query)
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
	return mdbNsInt.Mutate(ctx, mus)
}

func IntQuery(ctx context.Context, query string) (*Response, error) {
	resp, err := mdbNsInt.Query(ctx, query)
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
