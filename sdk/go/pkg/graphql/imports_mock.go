//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package graphql

import "github.com/hypermodeAI/functions-go/pkg/testutils"

var GraphqlQueryCallStack = testutils.NewCallStack()

func hostExecuteGQL(hostName, statement, variables *string) *string {
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
