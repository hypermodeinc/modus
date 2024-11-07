//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package modusdb

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var ModusDBQueryCallStack = testutils.NewCallStack()
var ModusDBAlterSchemaCallStack = testutils.NewCallStack()
var ModusDBDropAttrCallStack = testutils.NewCallStack()
var ModusDBDropAllCallStack = testutils.NewCallStack()
var ModusDBMutateCallStack = testutils.NewCallStack()

func hostAlterSchema(schema *string) *string {
	ModusDBAlterSchemaCallStack.Push(schema)

	success := "Success"

	return &success
}

func hostDropData() *string {
	ModusDBDropAttrCallStack.Push()

	success := "Success"

	return &success
}

func hostDropAll() *string {
	ModusDBDropAllCallStack.Push()

	success := "Success"

	return &success
}

func hostQuery(query *string) *Response {
	ModusDBQueryCallStack.Push(query)

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

func hostMutate(mutationReq *MutationRequest) *map[string]uint64 {
	ModusDBMutateCallStack.Push(mutationReq)

	return &map[string]uint64{
		"uid1": 1,
		"uid2": 2,
	}
}
