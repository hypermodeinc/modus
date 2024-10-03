//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db_test

import (
	"testing"

	"github.com/hypermodeAI/functions-go/pkg/db"
	"github.com/hypermodeAI/functions-go/pkg/utils"
)

var (
	testHostName = "mydb"
	testDbType   = "somedbtype"
)

func TestExecute(t *testing.T) {
	affected, err := db.Execute(testHostName, testDbType, db.MockExecuteStatement, db.MockExecuteParameters...)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if affected != 3 {
		t.Errorf("Expected 3 rows affected, but received: %d", affected)
	}

	testCallStack(t, db.MockExecuteStatement, db.MockExecuteParameters)
}

func TestQuery(t *testing.T) {
	rows, affected, err := db.Query[map[string]any](testHostName, testDbType, db.MockQueryStatement, db.MockQueryParameters...)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if affected != 3 {
		t.Errorf("Expected 3 rows affected, but received: %d", affected)
	}
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows, but received: %d", len(rows))
	}

	for i, row := range rows {
		n, err := utils.ConvertInterfaceTo[int](row["id"])
		if err != nil {
			t.Fatalf("Expected no error, but received: %s", err)
		}
		if n != i+1 {
			t.Errorf("Expected id: %d, but received: %v", i+1, row["id"])
		}
		switch i {

		}
	}

	testCallStack(t, db.MockQueryStatement, db.MockQueryParameters)
}

func TestQueryScalar(t *testing.T) {
	result, affected, err := db.QueryScalar[int](testHostName, testDbType, db.MockQueryScalarStatement, db.MockQueryScalarParameters...)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if affected != 1 {
		t.Errorf("Expected 1 rows affected, but received: %d", affected)
	}

	if result != 3 {
		t.Errorf("Expected result: 3, but received: %d", result)
	}

	testCallStack(t, db.MockQueryScalarStatement, db.MockQueryScalarParameters)
}

func testCallStack(t *testing.T, expectedStatement string, expectedParams []any) {
	values := db.DatabaseQueryCallStack.Pop()
	if values == nil {
		t.Error("Expected a query, but none was found.")
	} else {
		if len(values) != 4 {
			t.Errorf("Expected 4 values, but received: %d", len(values))
		}
		receivedHostName := values[0].(*string)
		if *receivedHostName != testHostName {
			t.Errorf("Expected hostName: \"%s\", but received: \"%s\"", testHostName, *receivedHostName)
		}
		receivedDbType := values[1].(*string)
		if *receivedDbType != testDbType {
			t.Errorf("Expected dbType: \"%s\", but received: \"%s\"", testDbType, *receivedDbType)
		}
		receivedStatement := values[2].(*string)
		if *receivedStatement != expectedStatement {
			t.Errorf("Expected statement: \"%s\", but received: \"%s\"", expectedStatement, *receivedStatement)
		}
		receivedJson := values[3].(*string)

		bytes, _ := utils.JsonSerialize(expectedParams)
		expectedParamsJson := string(bytes)
		if *receivedJson != expectedParamsJson {
			t.Errorf("Expected paramsJson: \"%s\", but received: \"%s\"", expectedParamsJson, *receivedJson)
		}
	}
}
