/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"slices"
	"testing"
)

func TestArrayInput_i32(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []int32{1, 2, 3}

	_, err := f.InvokeFunction("testArrayInput_i32", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []string{"abc", "def", "ghi"}

	_, err := f.InvokeFunction("testArrayInput_string", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput_i32(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput_i32")
	if err != nil {
		t.Fatal(err)
	}

	arr := []int32{1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int32); !ok {
		t.Errorf("expected a []int32, got %T", result)
	} else if !slices.Equal(arr, r) {
		t.Errorf("expected %x, got %x", arr, r)
	}
}

func TestArrayOutput_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput_string")
	if err != nil {
		t.Fatal(err)
	}

	arr := []string{"abc", "def", "ghi"}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]string); !ok {
		t.Errorf("expected a []string, got %T", result)
	} else if !slices.Equal(arr, r) {
		t.Errorf("expected %x, got %x", arr, r)
	}
}

type TestObject1 struct {
	A int32
	B int32
}

func TestArrayIteration(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	// Note, the below works, but we can make it fail easily
	// if we make the array larger than the elements we provide.
	// That's because the test function uses a non-nullable type for
	// the array elements and thus it will fail if it encounters an
	// undefined value.

	arr := make([]*TestObject1, 3)
	arr[0] = &TestObject1{A: 1, B: 2}
	arr[1] = &TestObject1{A: 3, B: 4}
	arr[2] = &TestObject1{A: 5, B: 6}

	_, err := f.InvokeFunction("testArrayIteration", arr)
	if err != nil {
		t.Fatal(err)
	}
}
