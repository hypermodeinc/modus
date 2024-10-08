/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang_test

import (
	"bytes"
	"reflect"
	"slices"
	"testing"
)

func TestSliceInput_byte(t *testing.T) {
	fnName := "testSliceInput_byte"
	s := []byte{1, 2, 3, 4}

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func TestSliceInput_intPtr(t *testing.T) {
	fnName := "testSliceInput_intPtr"
	s := getIntPtrSlice()

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func TestSliceInput_string(t *testing.T) {
	fnName := "testSliceInput_string"
	s := []string{"abc", "def", "ghi"}

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func TestSliceInput_stringPtr(t *testing.T) {
	fnName := "testSliceInput_stringPtr"
	s := getStringPtrSlice()

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func TestSliceOutput_byte(t *testing.T) {
	fnName := "testSliceOutput_byte"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{0x01, 0x02, 0x03, 0x04}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]byte); !ok {
		t.Errorf("expected a []byte, got %T", result)
	} else if !bytes.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestSliceOutput_intPtr(t *testing.T) {
	fnName := "testSliceOutput_intPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getIntPtrSlice()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]*int); !ok {
		t.Errorf("expected a []*int, got %T", result)
	} else if !slices.EqualFunc(expected, r, func(a, b *int) bool { return *a == *b }) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestSliceOutput_string(t *testing.T) {
	fnName := "testSliceOutput_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"abc", "def", "ghi"}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]string); !ok {
		t.Errorf("expected a []string, got %T", result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestSliceOutput_stringPtr(t *testing.T) {
	fnName := "testSliceOutput_stringPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getStringPtrSlice()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]*string); !ok {
		t.Errorf("expected a []*string, got %T", result)
	} else if !slices.EqualFunc(expected, r, func(a, b *string) bool { return *a == *b }) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func getIntPtrSlice() []*int {
	a := 11
	b := 22
	c := 33
	return []*int{&a, &b, &c}
}

func getStringPtrSlice() []*string {
	a := "abc"
	b := "def"
	c := "ghi"
	return []*string{&a, &b, &c}
}

func TestSliceInput_string_nil(t *testing.T) {
	fnName := "testSliceInput_string_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestSliceOutput_string_nil(t *testing.T) {
	fnName := "testSliceOutput_string_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestSliceInput_string_empty(t *testing.T) {
	fnName := "testSliceInput_string_empty"
	s := []string{}

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func TestSliceOutput_string_empty(t *testing.T) {
	fnName := "testSliceOutput_string_empty"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]string); !ok {
		t.Errorf("expected a []string, got %T", result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestSliceInput_int32_empty(t *testing.T) {
	fnName := "testSliceInput_int32_empty"
	s := []int32{}

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func TestSliceOutput_int32_empty(t *testing.T) {
	fnName := "testSliceOutput_int32_empty"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []int32{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int32); !ok {
		t.Errorf("expected a []int32, got %T", result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func Test2DSliceInput_string(t *testing.T) {
	fnName := "test2DSliceInput_string"
	s := [][]string{
		{"abc", "def", "ghi"},
		{"jkl", "mno", "pqr"},
		{"stu", "vwx", "yz"},
	}

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func Test2DSliceOutput_string(t *testing.T) {
	fnName := "test2DSliceOutput_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]string{
		{"abc", "def", "ghi"},
		{"jkl", "mno", "pqr"},
		{"stu", "vwx", "yz"},
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([][]string); !ok {
		t.Errorf("expected a []string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func Test2DSliceInput_string_nil(t *testing.T) {
	fnName := "test2DSliceInput_string_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func Test2DSliceOutput_string_nil(t *testing.T) {
	fnName := "test2DSliceOutput_string_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func Test2DSliceInput_string_empty(t *testing.T) {
	fnName := "test2DSliceInput_string_empty"
	s := [][]string{}

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func Test2DSliceOutput_string_empty(t *testing.T) {
	fnName := "test2DSliceOutput_string_empty"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]string{}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([][]string); !ok {
		t.Errorf("expected a []string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func Test2DSliceInput_string_innerNil(t *testing.T) {
	fnName := "test2DSliceInput_string_innerNil"
	s := [][]string{nil}

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
}

func Test2DSliceOutput_string_innerNil(t *testing.T) {
	fnName := "test2DSliceOutput_string_innerNil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]string{nil}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([][]string); !ok {
		t.Errorf("expected a []string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}
