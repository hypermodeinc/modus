//go:build wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import "unsafe"

//go:noescape
//go:wasmimport modus_test add
func hostAdd(a, b int32) int32

//go:noescape
//go:wasmimport modus_test echo1
func hostEcho1(message *string) *string

//go:noescape
//go:wasmimport modus_test echo2
func hostEcho2(message *string) *string

//go:noescape
//go:wasmimport modus_test echo3
func hostEcho3(message *string) *string

//go:noescape
//go:wasmimport modus_test echo4
func hostEcho4(message *string) *string

//go:noescape
//go:wasmimport modus_test encodeStrings1
func _hostEncodeStrings1(items unsafe.Pointer) *string

//modus:import modus_test encodeStrings1
func hostEncodeStrings1(items *[]string) *string {
	return _hostEncodeStrings1(unsafe.Pointer(items))
}

//go:noescape
//go:wasmimport modus_test encodeStrings2
func _hostEncodeStrings2(items unsafe.Pointer) *string

//modus:import modus_test encodeStrings2
func hostEncodeStrings2(items *[]*string) *string {
	return _hostEncodeStrings2(unsafe.Pointer(items))
}
