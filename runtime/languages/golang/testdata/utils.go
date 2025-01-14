/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"

	"github.com/hypermodeinc/modus/sdk/go/pkg/console"
)

func fail(msg string) {
	console.Error(msg)
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

// HasNil returns true if the given interface value is nil, or contains a pointer, interface, slice, or map that is nil.
func hasNil(x any) bool {
	if x == nil {
		return true
	}

	rv := reflect.ValueOf(x)
	if canBeNil(rv.Type()) {
		return rv.IsNil()
	}

	return false
}

func canBeNil(rt reflect.Type) bool {
	switch rt.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}
