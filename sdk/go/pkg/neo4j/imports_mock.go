//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4j

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var Neo4jQueryCallStack = testutils.NewCallStack()

func hostExecuteQuery(connection, dbName, query, parameters *string) *EagerResult {
	Neo4jQueryCallStack.Push(connection, dbName, query, parameters)

	keys := []string{"key1", "key2"}
	values := []string{"value1", "value2"}
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
