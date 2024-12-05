//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4j

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var DgraphQueryCallStack = testutils.NewCallStack()
var DgraphAlterSchemaCallStack = testutils.NewCallStack()
var DgraphDropAttrCallStack = testutils.NewCallStack()
var DgraphDropAllCallStack = testutils.NewCallStack()

func hostExecuteQuery(hostName, dbName, query, parameters *string) *EagerResult {
	DgraphQueryCallStack.Push(hostName, dbName, query, parameters)

	keys := []string{"key1", "key2"}
	values := []any{"value1", "value2"}
	record := &Record{
		Keys:   keys,
		Values: values,
	}
	records := []*Record{record}

	return &EagerResult{
		Keys:    keys,
		Records: records,
	}
}
