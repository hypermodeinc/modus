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
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Represents a Dgraph request.
type Request struct {
	Query     *Query
	Mutations []*Mutation
}

// Represents a Dgraph query.
type Query struct {
	Query     string
	Variables map[string]string
}

// Represents a Dgraph mutation.
type Mutation struct {

	// JSON for setting data
	SetJson string

	// JSON for deleting data
	DelJson string

	// RDF N-Quads for setting data
	SetNquads string

	// RDF N-Quads for deleting data
	DelNquads string

	// Condition for the mutation, as a DQL @if directive
	Condition string
}

// Represents a Dgraph response.
type Response struct {
	Json string
	Uids map[string]string
}

// Creates a new Dgraph query.
func NewQuery(query string) *Query {
	return &Query{
		Query:     query,
		Variables: make(map[string]string),
	}
}

// Adds a variable to the query.
func (q *Query) WithVariable(key string, value any) *Query {
	switch value := value.(type) {
	case string:
		// strings are passed as-is, without extra encoding
		q.Variables[key] = value
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, bool,
		[]float32:
		s, _ := json.Marshal(value)
		q.Variables[key] = string(s)
	default:
		panic(fmt.Errorf("unsupported DQL variable type: %T (must be string, integer, float, bool, or []float32)", value))
	}

	return q
}

// Creates a new Dgraph mutation.
func NewMutation() *Mutation {
	return &Mutation{}
}

// Adds a JSON string for setting data to the mutation.
func (m *Mutation) WithSetJson(json string) *Mutation {
	m.SetJson = json
	return m
}

// Adds a JSON string for deleting data to the mutation.
func (m *Mutation) WithDelJson(json string) *Mutation {
	m.DelJson = json
	return m
}

// Adds RDF N-Quads for setting data to the mutation.
func (m *Mutation) WithSetNquads(nquads string) *Mutation {
	m.SetNquads = nquads
	return m
}

// Adds RDF N-Quads for deleting data to the mutation.
func (m *Mutation) WithDelNquads(nquads string) *Mutation {
	m.DelNquads = nquads
	return m
}

// Adds a condition to the mutation, as a DQL @if directive.
func (m *Mutation) WithCondition(cond string) *Mutation {
	m.Condition = cond
	return m
}

// Deprecated: Use ExecuteQuery or ExecuteMutations instead.
func Execute(connection string, request *Request) (*Response, error) {
	response := hostExecuteQuery(&connection, request)
	if response == nil {
		return nil, errors.New("Failed to execute the DQL query.")
	}

	return response, nil
}

// Executes a DQL query on the Dgraph database, optionally with mutations.
func ExecuteQuery(connection string, query *Query, mutations ...*Mutation) (*Response, error) {
	request := &Request{
		Query:     query,
		Mutations: mutations,
	}

	return Execute(connection, request)
}

// Executes one or more mutations on the Dgraph database.
func ExecuteMutations(connection string, mutations ...*Mutation) (*Response, error) {
	request := &Request{
		Mutations: mutations,
	}

	return Execute(connection, request)
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

// Ensures proper escaping of RDF string literals
func EscapeRDF(value string) string {
	var sb strings.Builder
	for _, r := range value {
		switch r {
		case '\\':
			sb.WriteString(`\\`)
		case '"':
			sb.WriteString(`\"`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		default:
			// handle control characters
			if r < 0x20 || (r >= 0x7F && r <= 0x9F) {
				sb.WriteString(fmt.Sprintf(`\u%04X`, r))
			} else {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}
