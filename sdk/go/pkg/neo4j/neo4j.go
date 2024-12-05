/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4j

import "github.com/hypermodeinc/modus/sdk/go/pkg/utils"

type Neo4jOption func(*neo4jOptions)

type neo4jOptions struct {
	dbName string
}

func WithDbName(dbName string) Neo4jOption {
	return func(o *neo4jOptions) {
		o.dbName = dbName
	}
}

type EagerResult struct {
	Keys    []string
	Records []*Record
}

type Record struct {
	Values []string
	Keys   []string
}

/**
 *
 * Executes a query or mutation on the Neo4j database.
 *
 * @param hostName - the name of the host
 * @param query - the query to execute
 * @param parameters - the parameters to pass to the query
 */
func ExecuteQuery(hostName, query string, parameters map[string]any, opts ...Neo4jOption) (*EagerResult, error) {
	dbOpts := &neo4jOptions{
		dbName: "neo4j",
	}

	for _, opt := range opts {
		opt(dbOpts)
	}

	bytes, err := utils.JsonSerialize(parameters)
	if err != nil {
		return nil, err
	}

	parametersJson := string(bytes)

	response := hostExecuteQuery(&hostName, &dbOpts.dbName, &query, &parametersJson)

	return response, nil
}
