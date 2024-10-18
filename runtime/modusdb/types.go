/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package modusdb

type Mutation struct {
	SetJson   string
	DelJson   string
	SetNquads string
	DelNquads string
	Condition string
}

type Response struct {
	Json    string
	Txn     *TxnContext
	Latency *Latency
	Metrics *Metrics
	Uids    map[string]string
	Rdf     []byte
	Hdrs    map[string]*ListOfString
}

type TxnContext struct {
	StartTs  uint64
	CommitTs uint64
	Aborted  bool
	Keys     []string
	Preds    []string
	Hash     string
}

type Latency struct {
	ParsingNs         uint64
	ProcessingNs      uint64
	EncodingNs        uint64
	AssignTimestampNs uint64
	TotalNs           uint64
}

type Metrics struct {
	NumUids map[string]uint64
}

type ListOfString struct {
	Value []string
}
