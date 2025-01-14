//go:build wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraph

import "unsafe"

//go:noescape
//go:wasmimport modus_dgraph_client executeQuery
func _hostExecuteQuery(hostName *string, request unsafe.Pointer) unsafe.Pointer

//modus:import modus_dgraph_client executeQuery
func hostExecuteQuery(hostName *string, request *Request) *Response {
	response := _hostExecuteQuery(hostName, unsafe.Pointer(request))
	if response == nil {
		return nil
	}
	return (*Response)(response)
}

//go:noescape
//go:wasmimport modus_dgraph_client alterSchema
func hostAlterSchema(hostName, schema *string) *string

//go:noescape
//go:wasmimport modus_dgraph_client dropAttribute
func hostDropAttribute(hostName, attr *string) *string

//go:noescape
//go:wasmimport modus_dgraph_client dropAllData
func hostDropAllData(hostName *string) *string
