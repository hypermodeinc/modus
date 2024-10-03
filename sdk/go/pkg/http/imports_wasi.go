//go:build wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package http

import "unsafe"

//go:noescape
//go:wasmimport hypermode httpFetch
func _fetch(request unsafe.Pointer) unsafe.Pointer

//hypermode:import hypermode httpFetch
func fetch(request *Request) *Response {
	response := _fetch(unsafe.Pointer(request))
	if response == nil {
		return nil
	}
	return (*Response)(response)
}
