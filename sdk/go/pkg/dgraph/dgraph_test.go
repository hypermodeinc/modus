//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraph_test

import (
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"
)

var (
	connection = "mydgraph"
	request    = &dgraph.Request{
		Query: &dgraph.Query{
			Query: "query",
		},
	}
	schema = "schema"
)

func TestExecuteDQL(t *testing.T) {
	response, err := dgraph.Execute(connection, request)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}
	if response == nil {
		t.Fatalf("Expected a response, but received nil")
		return
	}

	if response.Json != `{"data": {"query": "query"}}` {
		t.Errorf("Expected json: `{\"data\": {\"query\": \"query\"}}`, but received: %s", response.Json)
	}

	if response.Uids == nil {
		t.Fatalf("Expected a uids response, but received nil")
	}
	if len(response.Uids) != 2 {
		t.Errorf("Expected 2 uids, but received: %d", len(response.Uids))
	}
	if response.Uids["uid1"] != "0x1" {
		t.Errorf("Expected uid1: \"0x1\", but received: %s", response.Uids["uid1"])
	}
	if response.Uids["uid2"] != "0x2" {
		t.Errorf("Expected uid2: \"0x2\", but received: %s", response.Uids["uid2"])
	}

	expectedConnection := &connection
	expectedReq := &dgraph.Request{
		Query: &dgraph.Query{
			Query: "query",
		},
	}

	values := dgraph.DgraphQueryCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(expectedConnection, values[0]) {
			t.Errorf("Expected connection: %s, but received: %s", *expectedConnection, values[0])
		}
		if !reflect.DeepEqual(expectedReq, values[1]) {
			t.Errorf("Expected request: %v, but received: %v", expectedReq, values[1])
		}
	}
}

func TestAlterSchema(t *testing.T) {
	err := dgraph.AlterSchema(connection, schema)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}

	expectedConnection := &connection
	expectedSchema := &schema
	values := dgraph.DgraphAlterSchemaCallStack.Pop()
	if values == nil {
		t.Error("Expected a schema, but none was found.")
	} else {
		if !reflect.DeepEqual(expectedConnection, values[0]) {
			t.Errorf("Expected connection: %s, but received: %s", *expectedConnection, values[0])
		}
		if !reflect.DeepEqual(expectedSchema, values[1]) {
			t.Errorf("Expected schema: %v, but received: %v", *expectedSchema, values[1])
		}
	}

}

func TestDropAttr(t *testing.T) {
	attr := "attr"
	err := dgraph.DropAttr(connection, attr)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}

	expectedConnection := &connection
	expectedAttr := &attr

	values := dgraph.DgraphDropAttrCallStack.Pop()
	if values == nil {
		t.Error("Expected an attribute, but none was found.")
	} else {
		if !reflect.DeepEqual(expectedConnection, values[0]) {
			t.Errorf("Expected connection: %s, but received: %s", *expectedConnection, values[0])
		}
		if !reflect.DeepEqual(expectedAttr, values[1]) {
			t.Errorf("Expected attr: %v, but received: %v", *expectedAttr, values[1])
		}
	}
}

func TestDropAll(t *testing.T) {
	err := dgraph.DropAll(connection)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}

	expectedConnection := &connection

	values := dgraph.DgraphDropAllCallStack.Pop()
	if values == nil {
		t.Error("Expected a connection, but none was found.")
	} else {
		if !reflect.DeepEqual(expectedConnection, values[0]) {
			t.Errorf("Expected connection: %s, but received: %s", *expectedConnection, values[0])
		}
	}
}
