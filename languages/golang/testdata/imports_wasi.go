//go:build wasip1

package main

//go:noescape
//go:wasmimport test add
func hostAdd(a, b int) int

//go:noescape
//go:wasmimport test echo
func hostEcho(message *string) *string
