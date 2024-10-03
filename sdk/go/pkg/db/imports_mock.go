//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"github.com/hypermodeAI/functions-go/pkg/testutils"
)

var DatabaseQueryCallStack = testutils.NewCallStack()

var (
	MockExecuteStatement  = "UPDATE users SET name = $1 age = $2 WHERE id = $3"
	MockExecuteParameters = []any{"John", 30, 1}

	MockQueryStatement  = "SELECT * FROM users WHERE age >= $1 and age < $2 and active = $3"
	MockQueryParameters = []any{18, 21, true}

	MockQueryScalarStatement  = "SELECT COUNT(*) FROM users WHERE age >= $1 and age < $2 and active = $3"
	MockQueryScalarParameters = []any{0, 18, false}
)

func databaseQuery(hostName, dbType, statement, paramsJson *string) *HostQueryResponse {
	DatabaseQueryCallStack.Push(hostName, dbType, statement, paramsJson)

	switch *statement {
	case MockExecuteStatement:
		return &HostQueryResponse{
			Error:        nil,
			ResultJson:   nil,
			RowsAffected: 3,
		}
	case MockQueryStatement:
		result := `[{"id":1,"name":"Alice","age":18,"active":true},{"id":2,"name":"Bob","age":19,"active":true},{"id":3,"name":"Charlie","age":20,"active":true}]`
		return &HostQueryResponse{
			Error:        nil,
			ResultJson:   &result,
			RowsAffected: 3,
		}
	case MockQueryScalarStatement:
		result := `[{"count":3}]`
		return &HostQueryResponse{
			Error:        nil,
			ResultJson:   &result,
			RowsAffected: 1,
		}
	}

	panic("un-mocked database query")
}
