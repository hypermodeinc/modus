/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

func TestSliceInput_byte(val []byte) {
	expected := []byte{1, 2, 3, 4}
	assertSlicesEqual(expected, val)
}

func TestSliceOutput_byte() []byte {
	return []byte{1, 2, 3, 4}
}

func TestSliceInput_string(val []string) {
	expected := []string{"abc", "def", "ghi"}
	assertSlicesEqual(expected, val)
}

func TestSliceOutput_string() []string {
	return []string{"abc", "def", "ghi"}
}

func TestSliceInput_string_nil(val []string) {
	assertNil(val)
}

func TestSliceOutput_string_nil() []string {
	return nil
}

func TestSliceInput_string_empty(val []string) {
	expected := []string{}
	assertSlicesEqual(expected, val)
}

func TestSliceOutput_string_empty() []string {
	return []string{}
}

func TestSliceInput_int32_empty(val []int32) {
	expected := []int32{}
	assertSlicesEqual(expected, val)
}

func TestSliceOutput_int32_empty() []int32 {
	return []int32{}
}

func TestSliceInput_stringPtr(val []*string) {
	expected := getStringPtrSlice()
	assertPtrSlicesEqual(expected, val)
}

func TestSliceOutput_stringPtr() []*string {
	return getStringPtrSlice()
}

func getStringPtrSlice() []*string {
	a := "abc"
	b := "def"
	c := "ghi"
	return []*string{&a, &b, &c}
}

func TestSliceInput_intPtr(val []*int) {
	expected := getIntPtrSlice()
	assertPtrSlicesEqual(expected, val)
}

func TestSliceOutput_intPtr() []*int {
	return getIntPtrSlice()
}

func getIntPtrSlice() []*int {
	a := 11
	b := 22
	c := 33
	return []*int{&a, &b, &c}
}

func Test2DSliceInput_string(val [][]string) {
	expected := [][]string{
		{"abc", "def", "ghi"},
		{"jkl", "mno", "pqr"},
		{"stu", "vwx", "yz"},
	}
	assertEqual(expected, val)
}

func Test2DSliceOutput_string() [][]string {
	return [][]string{
		{"abc", "def", "ghi"},
		{"jkl", "mno", "pqr"},
		{"stu", "vwx", "yz"},
	}
}

func Test2DSliceInput_string_nil(val [][]string) {
	assertNil(val)
}

func Test2DSliceOutput_string_nil() [][]string {
	return nil
}

func Test2DSliceInput_string_empty(val [][]string) {
	expected := [][]string{}
	assertEqual(expected, val)
}

func Test2DSliceOutput_string_empty() [][]string {
	return [][]string{}
}

func Test2DSliceInput_string_innerNil(val [][]string) {
	expected := [][]string{nil}
	assertEqual(expected, val)
}

func Test2DSliceOutput_string_innerNil() [][]string {
	return [][]string{nil}
}
