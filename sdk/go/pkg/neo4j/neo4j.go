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
	"encoding/json"
	"fmt"
	"time"

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

type RecordValue interface {
	bool | int64 | float64 | string |
		time.Time |
		[]byte | []any | map[string]any |
		Node | Relationship | Path
}

type Entity interface {
	GetElementId() string
	GetProperties() map[string]any
}

type Node struct {
	ElementId string         `json:"ElementId"`
	Labels    []string       `json:"Labels"`
	Props     map[string]any `json:"Props"`
}

func (n *Node) GetElementId() string {
	return n.ElementId
}

func (n *Node) GetProperties() map[string]any {
	return n.Props
}

type Relationship struct {
	ElementId      string         `json:"ElementId"`
	StartElementId string         `json:"StartElementId"`
	EndElementId   string         `json:"EndElementId"`
	Type           string         `json:"Type"`
	Props          map[string]any `json:"Props"`
}

func (r *Relationship) GetElementId() string {
	return r.ElementId
}

func (r *Relationship) GetProperties() map[string]any {
	return r.Props
}

type Path struct {
	Nodes         []Node         `json:"Nodes"`
	Relationships []Relationship `json:"Relationships"`
}

type PropertyValue interface {
	bool | int64 | float64 | string |
		time.Time | []byte | []any
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

func GetRecordValue[T RecordValue](record *Record, key string) (T, error) {
	var val T
	for i, k := range record.Keys {
		if k == key {
			err := json.Unmarshal([]byte(record.Values[i]), &val)
			if err != nil {
				return *new(T), err
			} else {
				return val, nil
			}
		}
	}
	return *new(T), fmt.Errorf("Key not found in record")

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

func GetProperty[T PropertyValue](e Entity, key string) (T, error) {
	var val T
	rawVal, ok := e.GetProperties()[key]
	if !ok {
		return *new(T), fmt.Errorf("Key not found in node")
	}
	switch any(val).(type) {
	case int64:
		float64Val, ok := rawVal.(float64)
		if !ok {
			return *new(T), fmt.Errorf("expected value to have type int64 but found type %T", rawVal)
		}
		return any(int64(float64Val)).(T), nil
	default:
		val, ok = rawVal.(T)
		if !ok {
			zeroValue := *new(T)
			return zeroValue, fmt.Errorf("expected value to have type %T but found type %T", zeroValue, rawVal)
		}
		return val, nil
	}

}
