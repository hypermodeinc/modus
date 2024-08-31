//go:build wasip1

package main

//go:noescape
//go:wasmimport test add
func hostAdd(a, b int) int

//go:noescape
//go:wasmimport test echo1
func hostEcho1(message *string) *string

//go:noescape
//go:wasmimport test echo2
func hostEcho2(message *string) *string

//go:noescape
//go:wasmimport test echo3
func hostEcho3(message *string) *string

//go:noescape
//go:wasmimport test echo4
func hostEcho4(message *string) *string
