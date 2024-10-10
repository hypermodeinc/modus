//go:build wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package http

import "unsafe"

//go:noescape
//go:wasmimport modus_http_client httpFetch
func _fetch(request unsafe.Pointer) unsafe.Pointer

//hypermode:import hypermode httpFetch
func fetch(request *Request) *Response {
	response := _fetch(unsafe.Pointer(request))
	if response == nil {
		return nil
	}
	return (*Response)(response)
}
