//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphql_test

import (
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/sdk/go/pkg/graphql"
	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

func TestExecuteGQL(t *testing.T) {
	hostName := "test"
	statement := "query { test }"
	variables := map[string]any{
		"key":  "value",
		"key2": 2,
	}
	response, err := graphql.Execute[string](hostName, statement, variables)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}
	if response == nil {
		t.Fatalf("Expected a response, but received nil")
	}

	expectedHostName := &hostName
	expectedStatement := &statement
	bytes, err := utils.JsonSerialize(variables)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}

	expectedVariables := string(bytes)

	values := graphql.GraphqlQueryCallStack.Pop()
	if values == nil {
		t.Fatalf("Expected a call to hostExecuteGQL, but none was made")
	} else {
		if !reflect.DeepEqual(values[0], expectedHostName) {
			t.Errorf("Expected hostName: %s, but received: %s", *expectedHostName, values[0])
		}
		if !reflect.DeepEqual(values[1], expectedStatement) {
			t.Errorf("Expected statement: %s, but received: %s", *expectedStatement, values[1])
		}
		if !reflect.DeepEqual(*(values[2].(*string)), expectedVariables) {
			t.Errorf("Expected variables: %v, but received: %v", expectedVariables, *(values[2].(*string)))
		}
	}
}

func TestExecuteGQLWithErrors(t *testing.T) {
	hostName := "test"
	statement := "error"
	variables := map[string]any{
		"key":  "value",
		"key2": 2,
	}
	response, err := graphql.Execute[string](hostName, statement, variables)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}
	if response == nil {
		t.Fatalf("Expected a response, but received nil")
	}

	expectedHostName := &hostName
	expectedStatement := &statement
	bytes, err := utils.JsonSerialize(variables)
	if err != nil {
		t.Errorf("Expected no error, but received: %s", err.Error())
	}

	expectedVariables := string(bytes)

	values := graphql.GraphqlQueryCallStack.Pop()
	if values == nil {
		t.Fatalf("Expected a call to hostExecuteGQL, but none was made")
	} else {
		if !reflect.DeepEqual(values[0], expectedHostName) {
			t.Errorf("Expected hostName: %s, but received: %s", *expectedHostName, values[0])
		}
		if !reflect.DeepEqual(values[1], expectedStatement) {
			t.Errorf("Expected statement: %s, but received: %s", *expectedStatement, values[1])
		}
		if !reflect.DeepEqual(*(values[2].(*string)), expectedVariables) {
			t.Errorf("Expected variables: %v, but received: %v", expectedVariables, *(values[2].(*string)))
		}
	}
}
