package main

import (
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
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
	switch reflect.TypeOf(actual).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if reflect.ValueOf(actual).IsNil() {
			return
		}
	}
	fail(fmt.Sprintf("expected nil, got %v", actual))
}
