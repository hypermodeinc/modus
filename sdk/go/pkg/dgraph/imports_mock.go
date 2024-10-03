//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package dgraph

import "github.com/hypermodeAI/functions-go/pkg/testutils"

var DgraphQueryCallStack = testutils.NewCallStack()
var DgraphAlterSchemaCallStack = testutils.NewCallStack()
var DgraphDropAttrCallStack = testutils.NewCallStack()
var DgraphDropAllCallStack = testutils.NewCallStack()

func hostExecuteDQL(hostName *string, request *Request) *Response {
	DgraphQueryCallStack.Push(hostName, request)

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

func hostDgraphAlterSchema(hostName, schema *string) *string {
	DgraphAlterSchemaCallStack.Push(hostName, schema)

	success := "Success"

	return &success
}

func hostDgraphDropAttr(hostName, attr *string) *string {
	DgraphDropAttrCallStack.Push(hostName, attr)

	success := "Success"

	return &success
}

func hostDgraphDropAll(hostName *string) *string {
	DgraphDropAllCallStack.Push(hostName)

	success := "Success"

	return &success
}
