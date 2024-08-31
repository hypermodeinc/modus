package main

import (
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"unsafe"
)

func fail(msg string) {
	os.Stderr.WriteString(msg + "\n")
	os.Exit(1)
}

func assertEqual[T any](expect, actual T) {
	if !reflect.DeepEqual(expect, actual) {
		fail(fmt.Sprintf("expected %v, got %v", expect, actual))
	}
}

func assertSlicesEqual[T comparable](expect, actual []T) {
	if (expect == nil && actual != nil) || (expect != nil && actual == nil) || !slices.Equal(expect, actual) {
		fail(fmt.Sprintf("expected %v, got %v", expect, actual))
	}
}

func assertMapsEqual[T ~map[K]V, K, V comparable](expect, actual T) {
	if (expect == nil && actual != nil) || (expect != nil && actual == nil) || !maps.Equal(expect, actual) {
		fail(fmt.Sprintf("expected %v, got %v", expect, actual))
	}
}

func assertPtrSlicesEqual[T comparable](expect, actual []*T) {
	if (expect == nil && actual != nil) || (expect != nil && actual == nil) ||
		!slices.EqualFunc(expect, actual, ptrValsEqual) {
		fail(fmt.Sprintf("expected %v, got %v", expect, actual))
	}
}

func ptrValsEqual[T comparable](a, b *T) bool {
	return a == nil && b == nil || a != nil && b != nil && *a == *b
}

func assertNil[T any](actual T) {
	if !hasNil(actual) {
		fail(fmt.Sprintf("expected nil, got %v", actual))
	}
}

func getUnsafeDataPtr(x any) unsafe.Pointer {
	type iface struct {
		typ  unsafe.Pointer
		data unsafe.Pointer
	}

	internal := *(*iface)(unsafe.Pointer(&x))
	return internal.data
}

// hasNil returns true if the given interface value is nil, or contains a nil pointer.
func hasNil(x any) bool {
	return x == nil || uintptr(getUnsafeDataPtr(x)) == 0
}
