//go:build wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraph

import "unsafe"

//go:noescape
//go:wasmimport hypermode executeDQL
func _hostExecuteDQL(hostName *string, request unsafe.Pointer) unsafe.Pointer

//hypermode:import hypermode executeDQL
func hostExecuteDQL(hostName *string, request *Request) *Response {
	response := _hostExecuteDQL(hostName, unsafe.Pointer(request))
	if response == nil {
		return nil
	}
	return (*Response)(response)
}

//go:noescape
//go:wasmimport hypermode dgraphAlterSchema
func hostDgraphAlterSchema(hostName, schema *string) *string

//go:noescape
//go:wasmimport hypermode dgraphDropAttr
func hostDgraphDropAttr(hostName, attr *string) *string

//go:noescape
//go:wasmimport hypermode dgraphDropAll
func hostDgraphDropAll(hostName *string) *string
