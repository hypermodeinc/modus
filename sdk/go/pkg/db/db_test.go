//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db_test

import (
	"strings"
	"testing"

	"github.com/hypermodeinc/modus/sdk/go/pkg/db"
	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

var (
	testConnection = "mydb"
	testDbType     = "somedbtype"
)

func TestExecute(t *testing.T) {
	r, err := db.Execute(testConnection, testDbType, db.MockExecuteStatement, db.MockExecuteParameters...)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if r.RowsAffected != 3 {
		t.Errorf("Expected 3 rows affected, but received: %d", r.RowsAffected)
	}

	testCallStack(t, db.MockExecuteStatement, db.MockExecuteParameters, "exec")
}

func TestQuery(t *testing.T) {
	r, err := db.Query[map[string]any](testConnection, testDbType, db.MockQueryStatement, db.MockQueryParameters...)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if r.RowsAffected != 3 {
		t.Errorf("Expected 3 rows affected, but received: %d", r.RowsAffected)
	}
	if len(r.Rows) != 3 {
		t.Errorf("Expected 3 rows, but received: %d", len(r.Rows))
	}

	for i, row := range r.Rows {
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
	r, err := db.QueryScalar[int](testConnection, testDbType, db.MockQueryScalarStatement, db.MockQueryScalarParameters...)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if r.RowsAffected != 1 {
		t.Errorf("Expected 1 rows affected, but received: %d", r.RowsAffected)
	}

	if r.Value != 3 {
		t.Errorf("Expected result: 3, but received: %d", r.Value)
	}

	testCallStack(t, db.MockQueryScalarStatement, db.MockQueryScalarParameters)
}

func testCallStack(t *testing.T, expectedStatement string, expectedParams []any, flags ...string) {
	values := db.DatabaseQueryCallStack.Pop()
	if values == nil {
		t.Error("Expected a query, but none was found.")
	} else {
		if len(values) != 4 {
			t.Errorf("Expected 4 values, but received: %d", len(values))
		}
		receivedConnection := values[0].(*string)
		if *receivedConnection != testConnection {
			t.Errorf("Expected connection: \"%s\", but received: \"%s\"", testConnection, *receivedConnection)
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
		if len(flags) > 0 {
			expectedParamsJson = strings.Join(flags, ",") + ":" + expectedParamsJson
		}
		if *receivedJson != expectedParamsJson {
			t.Errorf("Expected paramsJson: \"%s\", but received: \"%s\"", expectedParamsJson, *receivedJson)
		}
	}
}
