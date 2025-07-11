//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4j

import "unsafe"

//go:noescape
//go:wasmimport modus_neo4j_client executeQuery
func _hostExecuteQuery(connection, dbName, query, parametersJson *string) unsafe.Pointer

//modus:import modus_neo4j_client executeQuery
func hostExecuteQuery(connection, dbName, query, parametersJson *string) *EagerResult {
	response := _hostExecuteQuery(connection, dbName, query, parametersJson)
	if response == nil {
		return nil
	}
	return (*EagerResult)(response)
}
