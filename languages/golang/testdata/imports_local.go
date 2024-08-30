//go:build !wasip1

package main

func hostAdd(a, b int) int {
	panic("stub")
}

func hostEcho(message *string) *string {
	panic("stub")
}
