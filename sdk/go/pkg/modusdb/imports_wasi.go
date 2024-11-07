//go:build wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package modusdb

import "unsafe"

//go:noescape
//go:wasmimport modus_modusdb_client dropAll
func hostDropAll() *string

//go:noescape
//go:wasmimport modus_modusdb_client dropData
func hostDropData() *string

//go:noescape
//go:wasmimport modus_modusdb_client alterSchema
func hostAlterSchema(schema *string) *string

//go:noescape
//go:wasmimport modus_modusdb_client mutate
func _hostMutate(mutationReq unsafe.Pointer) unsafe.Pointer

//modus:import modus_modusdb_client mutate
func hostMutate(mutationReq *MutationRequest) *map[string]uint64 {
	resp := _hostMutate(unsafe.Pointer(mutationReq))
	if resp == nil {
		return nil
	}
	return (*map[string]uint64)(resp)
}

//go:noescape
//go:wasmimport modus_modusdb_client query
func _hostQuery(query *string) unsafe.Pointer

//modus:import modus_modusdb_client query
func hostQuery(query *string) *Response {
	response := _hostQuery(query)
	if response == nil {
		return nil
	}
	return (*Response)(response)
}
