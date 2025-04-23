//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4j_test

import (
	"testing"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

var (
	connection = "myneo4j"
)

func TestRecordJSONMarshal(t *testing.T) {
	record := neo4j.Record{
		Keys:   []string{"name", "age", "friends"},
		Values: []string{`"Alice"`, `"42"`, `["Bob","Peter","Anna"]`},
	}

	// Marshal the record to JSON
	jsonBytes, err := record.JSONMarshal()
	if err != nil {
		t.Fatalf("Expected no error, but received: %v", err)
	}

	jsonStr := string(jsonBytes)
	expectedJSONStr := `{"name":"Alice","age":"42","friends":["Bob","Peter","Anna"]}`
	if jsonStr != expectedJSONStr {
		t.Errorf("Expected JSON: %s, but received: %s", expectedJSONStr, jsonStr)
	}
}

func TestExecuteQuery(t *testing.T) {
	dbName := "mydb"
	query := "query"
	parameters := map[string]any{
		"param1": "value1",
		"param2": "value2",
	}

	response, err := neo4j.ExecuteQuery(connection, query, parameters, neo4j.WithDbName(dbName))
	if err != nil {
		t.Fatalf("Expected no error, but received: %v", err)
	}
	if response == nil {
		t.Fatalf("Expected a response, but received nil")
		return
	}

	if len(response.Keys) != 2 {
		t.Errorf("Expected 2 keys, but received: %d", len(response.Keys))
	}
	if response.Keys[0] != "key1" {
		t.Errorf("Expected key1: \"key1\", but received: %s", response.Keys[0])
	}
	if response.Keys[1] != "key2" {
		t.Errorf("Expected key2: \"key2\", but received: %s", response.Keys[1])
	}

	if len(response.Records) != 1 {
		t.Errorf("Expected 1 record, but received: %d", len(response.Records))
	}
	if len(response.Records[0].Values) != 2 {
		t.Errorf("Expected 2 values, but received: %d", len(response.Records[0].Values))
	}
	if response.Records[0].Values[0] != "value1" {
		t.Errorf("Expected value1: \"value1\", but received: %s", response.Records[0].Values[0])
	}
	if response.Records[0].Values[1] != "value2" {
		t.Errorf("Expected value2: \"value2\", but received: %s", response.Records[0].Values[1])
	}
}
