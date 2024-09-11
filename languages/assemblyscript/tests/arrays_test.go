/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"slices"
	"testing"
)

func TestArrayInput_i8(t *testing.T) {
	arr := []int8{1, 2, 3}

	_, err := fixture.CallFunction(t, "testArrayInput_i8", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput_i8(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput_i8")
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
	arr := []int8{}

	_, err := fixture.CallFunction(t, "testArrayInput_i8_empty", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput_i8_empty(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput_i8_empty")
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
	_, err := fixture.CallFunction(t, "testArrayInput_i8_null", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput_i8_null(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput_i8_null")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Errorf("expected nil, got %T", result)
	}
}

func TestArrayInput_i32(t *testing.T) {
	arr := []int32{1, 2, 3}

	_, err := fixture.CallFunction(t, "testArrayInput_i32", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput_i32(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput_i32")
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
	arr := []float32{1, 2, 3}

	_, err := fixture.CallFunction(t, "testArrayInput_f32", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput_f32(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput_f32")
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
	arr := []string{"abc", "def", "ghi"}

	_, err := fixture.CallFunction(t, "testArrayInput_string", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput_string")
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

	_, err := fixture.CallFunction(t, "testArrayIteration", arr)
	if err != nil {
		t.Fatal(err)
	}
}
