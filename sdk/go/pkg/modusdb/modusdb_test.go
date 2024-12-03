//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package modusdb_test

import (
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/sdk/go/pkg/modusdb"
)

var (
	request = &modusdb.MutationRequest{
		Mutations: []*modusdb.Mutation{
			{
				SetJson:   "setJson",
				DelJson:   "delJson",
				SetNquads: "setNquads",
				DelNquads: "delNquads",
				Condition: "condition",
			},
		},
	}
	schema = "schema"
	query  = "query"
)

func TestDropData(t *testing.T) {
	err := modusdb.DropData()
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}

	values := modusdb.ModusDBDropAttrCallStack.Pop()
	if values != nil {
		t.Error("Expected no value, but was found: ", values)
	}
}

func TestAlterSchema(t *testing.T) {
	err := modusdb.AlterSchema(schema)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}

	expectedSchema := &schema
	values := modusdb.ModusDBAlterSchemaCallStack.Pop()
	if values == nil {
		t.Error("Expected a schema, but none was found.")
	} else {
		if !reflect.DeepEqual(expectedSchema, values[0]) {
			t.Errorf("Expected schema: %v, but received: %v", *expectedSchema, values[0])
		}
	}
}

func TestQuery(t *testing.T) {
	response, err := modusdb.Query(&query)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}
	if response == nil {
		t.Error("Expected a response, but none was found.")
	}

	expectedQuery := query
	values := modusdb.ModusDBQueryCallStack.Pop()
	if values == nil {
		t.Error("Expected a query, but none was found.")
	} else {
		if !reflect.DeepEqual(expectedQuery, values[0]) {
			t.Errorf("Expected query: %v, but received: %v", expectedQuery, values[0])
		}
	}
}

func TestMutate(t *testing.T) {
	response, err := modusdb.Mutate(request)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}
	if response == nil {
		t.Error("Expected a response, but none was found.")
	}

	expectedRequest := request
	values := modusdb.ModusDBMutateCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(expectedRequest, values[0]) {
			t.Errorf("Expected request: %v, but received: %v", expectedRequest, values[0])
		}
	}
}
