/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraph

import (
	"errors"
)

type Request struct {
	Query     *Query
	Mutations []*Mutation
}

type Query struct {
	Query     string
	Variables map[string]string
}

type Mutation struct {
	SetJson   string
	DelJson   string
	SetNquads string
	DelNquads string
	Condition string
}

type Response struct {
	Json string
	Uids map[string]string
}

/**
 *
 * Executes a DQL query or mutation on the Dgraph database.
 *
 * @param hostName - the name of the host
 * @param query - the query to execute
 * @param mutations - the mutations to execute
 * @returns The response from the Dgraph server
 */
func Execute(hostName string, request *Request) (*Response, error) {
	response := hostExecuteDQL(&hostName, request)
	if response == nil {
		return nil, errors.New("Failed to execute the DQL query.")
	}

	return response, nil
}

/**
 *
 * Alters the schema of the dgraph database
 *
 * @param hostName - the name of the host
 * @param schema - the schema to alter
 * @returns The response from the Dgraph server
 */
func AlterSchema(hostName, schema string) error {
	resp := hostDgraphAlterSchema(&hostName, &schema)
	if resp == nil {
		return errors.New("Failed to alter the schema.")
	}

	return nil
}

/**
 *
 * Drops an attribute from the schema.
 *
 * @param hostName - the name of the host
 * @param attr - the attribute to drop
 * @returns The response from the Dgraph server
 */
func DropAttr(hostName, attr string) error {
	response := hostDgraphDropAttr(&hostName, &attr)
	if response == nil {
		return errors.New("Failed to drop the attribute.")
	}

	return nil
}

/**
 *
 * Drops all data from the database.
 *
 * @param hostName - the name of the host
 * @returns The response from the Dgraph server
 */
func DropAll(hostName string) error {
	response := hostDgraphDropAll(&hostName)
	if response == nil {
		return errors.New("Failed to drop all data.")
	}

	return nil
}
