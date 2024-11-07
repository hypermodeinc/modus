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
	"os"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modusdb"
)

var mdb *modusdb.DB
var dataDir = "data"

func Init(ctx context.Context) {
	var err error
	// check if data directory exists
	shouldAlterSchema := false
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		shouldAlterSchema = true
	}
	mdb, err = modusdb.New(modusdb.NewDefaultConfig().WithDataDir("data"))
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb")
	}
	if shouldAlterSchema {
		err = mdb.AlterSchema(ctx, internalSchema)
		if err != nil {
			logger.Fatal(ctx).Err(err).Msg("Error initializing modusdb schema")
		}
	}
}

func Close(ctx context.Context) {
	if mdb != nil {
		mdb.Close()
	}
}

func DropAll(ctx context.Context) (string, error) {
	err := mdb.DropAll(ctx)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func DropData(ctx context.Context) (string, error) {
	err := mdb.DropData(ctx)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func AlterSchema(ctx context.Context, schema string) (string, error) {
	totalSchema := fmt.Sprintf("%s\n%s", internalSchema, schema)
	err := mdb.AlterSchema(ctx, totalSchema)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func Mutate(ctx context.Context, mutationReq MutationRequest) (map[string]uint64, error) {
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
	return mdb.Mutate(ctx, mus)
}

func Query(ctx context.Context, query string) (*Response, error) {
	resp, err := mdb.Query(ctx, query)
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
