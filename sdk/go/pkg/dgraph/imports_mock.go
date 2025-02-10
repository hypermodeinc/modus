//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraph

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var DgraphQueryCallStack = testutils.NewCallStack()
var DgraphAlterSchemaCallStack = testutils.NewCallStack()
var DgraphDropAttrCallStack = testutils.NewCallStack()
var DgraphDropAllCallStack = testutils.NewCallStack()

func hostExecuteQuery(connection *string, request *Request) *Response {
	DgraphQueryCallStack.Push(connection, request)

	json := `{"data": {"query": "query"}}`

	uids := map[string]string{
		"uid1": "0x1",
		"uid2": "0x2",
	}

	return &Response{
		Json: json,
		Uids: uids,
	}
}

func hostAlterSchema(connection, schema *string) *string {
	DgraphAlterSchemaCallStack.Push(connection, schema)

	success := "Success"

	return &success
}

func hostDropAttribute(connection, attr *string) *string {
	DgraphDropAttrCallStack.Push(connection, attr)

	success := "Success"

	return &success
}

func hostDropAllData(connection *string) *string {
	DgraphDropAllCallStack.Push(connection)

	success := "Success"

	return &success
}
