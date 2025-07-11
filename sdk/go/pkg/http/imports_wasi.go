//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package http

import "unsafe"

//go:noescape
//go:wasmimport modus_http_client fetch
func _hostFetch(request unsafe.Pointer) unsafe.Pointer

//modus:import modus_http_client fetch
func hostFetch(request *Request) *Response {
	response := _hostFetch(unsafe.Pointer(request))
	if response == nil {
		return nil
	}
	return (*Response)(response)
}
