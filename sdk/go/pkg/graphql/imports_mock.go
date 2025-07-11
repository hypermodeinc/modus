//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphql

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var GraphqlQueryCallStack = testutils.NewCallStack()

func hostExecuteQuery(connection, statement, variables *string) *string {
	GraphqlQueryCallStack.Push(connection, statement, variables)

	var json string
	if *statement == "error" {
		json = `{
			"errors": [
				{
					"message": "mock error message",
					"locations": [
						{
							"line": 1,
							"column": 2
						}
					],
					"path": ["mock", "path"]
				}
			],
			"data": "mock data"
		}`
	} else {
		json = `{
			"data": "Successful query result"
		}`
	}

	return &json
}
