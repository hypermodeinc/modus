/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"maps"
	"slices"
)

func TestMapInput_string_string(m map[string]string) {

	expected := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	assertMapsEqual(expected, m)
}

func TestMapPtrInput_string_string(m *map[string]string) {
	expected := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	assertMapsEqual(expected, *m)
}

func TestMapOutput_string_string() map[string]string {
	return map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}
}

func TestMapPtrOutput_string_string() *map[string]string {
	return &map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}
}

func TestIterateMap_string_string(m map[string]string) {
	println("len(m):", len(m))

	keys := slices.Sorted(maps.Keys(m))
	for _, k := range keys {
		println(k, m[k])
	}
}

func TestMapLookup_string_string(m map[string]string, key string) string {
	return m[key]
}

type TestStructWithMap struct {
	M map[string]string
}

func TestStructContainingMapInput_string_string(s TestStructWithMap) {
	expected := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	assertMapsEqual(expected, s.M)
}

func TestStructContainingMapOutput_string_string() TestStructWithMap {
	return TestStructWithMap{
		M: map[string]string{
			"a": "1",
			"b": "2",
			"c": "3",
		},
	}
}

func TestMapInput_int_float32(m map[int]float32) {
	expected := map[int]float32{
		1: 1.1,
		2: 2.2,
		3: 3.3,
	}

	assertMapsEqual(expected, m)
}

func TestMapOutput_int_float32() map[int]float32 {
	return map[int]float32{
		1: 1.1,
		2: 2.2,
		3: 3.3,
	}
}
