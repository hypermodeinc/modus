/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4j

import (
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

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

type Node struct {
	Id        int               `json:"Id"`
	ElementId string            `json:"ElementId"`
	Labels    []string          `json:"Labels"`
	Props     map[string]string `json:"Props"`
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

func GetRecordValue[T any](record *Record, key string) (T, error) {
	var val T
	for i, k := range record.Keys {
		if k == key {
			err := utils.JsonDeserialize([]byte(record.Values[i]), &val)
			if err != nil {
				return val, err
			} else {
				return val, nil
			}
		}
	}
	return val, fmt.Errorf("Key not found in record")
}

func (r *Record) Get(key string) (string, bool) {
	for i, k := range r.Keys {
		if k == key {
			return r.Values[i], true
		}
	}
	return "", false
}

func (r *Record) AsMap() map[string]string {
	result := make(map[string]string)
	for i, k := range r.Keys {
		result[k] = r.Values[i]
	}
	return result
}
