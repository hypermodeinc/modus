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
	"errors"
)

type MutationRequest struct {
	Mutations []*Mutation
}

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
	Rdf     []uint8
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

/**
 *
 * drops all data from the database
 *
 * @returns The response from the Dgraph server
 */
func DropData() error {
	response := hostDropData()
	if response == nil {
		return errors.New("Failed to drop all data.")
	}

	return nil
}

/**
 *
 * Alters the schema of the dgraph database
 *
 * @param hostName - the name of the host
 * @param schema - the schema to alter
 * @returns The response from the Dgraph server
 */
func AlterSchema(schema string) error {
	resp := hostAlterSchema(&schema)
	if resp == nil {
		return errors.New("Failed to alter the schema.")
	}

	return nil
}

/**
 *
 * Mutates the database
 *
 * @param mutationReq - the mutation request
 * @returns The response from the Dgraph server
 */
func Mutate(mutationReq *MutationRequest) (*map[string]uint64, error) {
	response := hostMutate(mutationReq)
	if response == nil {
		return nil, errors.New("Failed to mutate the database.")
	}

	return response, nil
}

/**
 *
 * Queries the database
 *
 * @param query - the query to execute
 * @returns The response from the Dgraph server
 */
func Query(query *string) (*Response, error) {
	response := hostQuery(query)
	if response == nil {
		return nil, errors.New("Failed to query the database.")
	}

	return response, nil
}
