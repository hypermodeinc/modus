//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraph

import "unsafe"

//go:noescape
//go:wasmimport modus_dgraph_client executeQuery
func _hostExecuteQuery(connection *string, request unsafe.Pointer) unsafe.Pointer

//modus:import modus_dgraph_client executeQuery
func hostExecuteQuery(connection *string, request *Request) *Response {
	response := _hostExecuteQuery(connection, unsafe.Pointer(request))
	if response == nil {
		return nil
	}
	return (*Response)(response)
}

//go:noescape
//go:wasmimport modus_dgraph_client alterSchema
func hostAlterSchema(connection, schema *string) *string

//go:noescape
//go:wasmimport modus_dgraph_client dropAttribute
func hostDropAttribute(connection, attr *string) *string

//go:noescape
//go:wasmimport modus_dgraph_client dropAllData
func hostDropAllData(connection *string) *string
