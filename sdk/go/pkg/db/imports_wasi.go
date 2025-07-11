//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import "unsafe"

//go:noescape
//go:wasmimport modus_sql_client executeQuery
func _hostExecuteQuery(connection, dbType, statement, paramsJson *string) unsafe.Pointer

//modus:import modus_sql_client executeQuery
func hostExecuteQuery(connection, dbType, statement, paramsJson *string) *HostQueryResponse {
	response := _hostExecuteQuery(connection, dbType, statement, paramsJson)
	if response == nil {
		return nil
	}
	return (*HostQueryResponse)(response)
}
