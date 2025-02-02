/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
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

// Executes a DQL query or mutation on the Dgraph database.
func Execute(connection string, request *Request) (*Response, error) {
	response := hostExecuteQuery(&connection, request)
	if response == nil {
		return nil, errors.New("Failed to execute the DQL query.")
	}

	return response, nil
}

// Alters the schema of the dgraph database
func AlterSchema(connection, schema string) error {
	resp := hostAlterSchema(&connection, &schema)
	if resp == nil {
		return errors.New("Failed to alter the schema.")
	}

	return nil
}

// Drops an attribute from the schema.
func DropAttr(connection, attr string) error {
	response := hostDropAttribute(&connection, &attr)
	if response == nil {
		return errors.New("Failed to drop the attribute.")
	}

	return nil
}

// Drops all data from the database.
func DropAll(connection string) error {
	response := hostDropAllData(&connection)
	if response == nil {
		return errors.New("Failed to drop all data.")
	}

	return nil
}
