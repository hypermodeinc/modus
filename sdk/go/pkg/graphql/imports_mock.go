//go:build !wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphql

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var GraphqlQueryCallStack = testutils.NewCallStack()

func hostExecuteQuery(hostName, statement, variables *string) *string {
	GraphqlQueryCallStack.Push(hostName, statement, variables)

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
