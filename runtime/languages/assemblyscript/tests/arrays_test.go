/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

func TestArrayInput_i8(t *testing.T) {
	fnName := "testArrayInput_i8"
	arr := []int8{1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput_i8(t *testing.T) {
	fnName := "testArrayOutput_i8"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []int8{1, 2, 3}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput_i8_empty(t *testing.T) {
	fnName := "testArrayInput_i8_empty"
	arr := []int8{}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput_i8_empty(t *testing.T) {
	fnName := "testArrayOutput_i8_empty"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []int8{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput_i8_null(t *testing.T) {
	fnName := "testArrayInput_i8_null"
	var arr []int8 = nil

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput_i8_null(t *testing.T) {
	fnName := "testArrayOutput_i8_null"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Errorf("expected nil, got %T", result)
	}
}

func TestArrayInput_i32(t *testing.T) {
	fnName := "testArrayInput_i32"
	arr := []int32{1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput_i32(t *testing.T) {
	fnName := "testArrayOutput_i32"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []int32{1, 2, 3}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput_f32(t *testing.T) {
	fnName := "testArrayInput_f32"
	arr := []float32{1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput_f32(t *testing.T) {
	fnName := "testArrayOutput_f32"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []float32{1, 2, 3}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]float32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput_string(t *testing.T) {
	fnName := "testArrayInput_string"
	arr := []string{"abc", "def", "ghi"}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput_string(t *testing.T) {
	fnName := "testArrayOutput_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"abc", "def", "ghi"}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput_string_2d(t *testing.T) {
	fnName := "testArrayInput_string_2d"
	arr := [][]string{
		{"abc", "def", "ghi"}, {"jkl", "mno", "pqr"}, {"stu", "vwx", "yz"},
	}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput_string_2d(t *testing.T) {
	fnName := "testArrayOutput_string_2d"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]string{
		{"abc", "def", "ghi"}, {"jkl", "mno", "pqr"}, {"stu", "vwx", "yz"},
	}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([][]string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput_string_2d_empty(t *testing.T) {
	fnName := "testArrayInput_string_2d_empty"
	arr := [][]string{}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput_string_2d_empty(t *testing.T) {
	fnName := "testArrayOutput_string_2d_empty"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]string{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([][]string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

type TestObject1 struct {
	A int32
	B int32
}

func TestArrayIteration(t *testing.T) {
	// Note, the below works, but we can make it fail easily
	// if we make the array larger than the elements we provide.
	// That's because the test function uses a non-nullable type for
	// the array elements and thus it will fail if it encounters an
	// undefined value.

	arr := make([]*TestObject1, 3)
	arr[0] = &TestObject1{A: 1, B: 2}
	arr[1] = &TestObject1{A: 3, B: 4}
	arr[2] = &TestObject1{A: 5, B: 6}

	fnName := "testArrayIteration"
	_, err := fixture.CallFunction(t, fnName, arr)
	if err != nil {
		t.Error(err)
	}
}
