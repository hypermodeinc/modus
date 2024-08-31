/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"maps"
	"reflect"
	"testing"
)

type TestStruct1 struct {
	A bool
}
type TestStruct2 struct {
	A bool
	B int
}

var testStruct1 = TestStruct1{
	A: true,
}

var testStruct2 = TestStruct2{
	A: true,
	B: 123,
}

var testStruct1AsMap = map[string]any{
	"A": true,
}

var testStruct2AsMap = map[string]any{
	"A": true,
	"B": 123,
}

func TestStructInput1(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput1", testStruct1); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput2(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput2", testStruct2); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput1(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput1", testStruct1); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput1", &testStruct1); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput2(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput2", testStruct2); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput2", &testStruct2); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput1_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput1", testStruct1AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput2_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput2", testStruct2AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput1_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput1", testStruct1AsMap); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput1", &testStruct1AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput2_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput2", testStruct2AsMap); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput2", &testStruct2AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput1_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput1_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput2_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput2_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestStructOutput1(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct1", reflect.TypeFor[TestStruct1]())

	result, err := f.InvokeFunction("testStructOutput1")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct1); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct1, r) {
		t.Errorf("expected %v, got %v", testStruct1, r)
	}
}

func TestStructOutput2(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct2", reflect.TypeFor[TestStruct2]())

	result, err := f.InvokeFunction("testStructOutput2")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct2); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct2, r) {
		t.Errorf("expected %v, got %v", testStruct2, r)
	}
}

func TestStructPtrOutput1(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct1", reflect.TypeFor[TestStruct1]())

	result, err := f.InvokeFunction("testStructPtrOutput1")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct1); !ok {
		t.Errorf("expected a pointer to a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct1, *r) {
		t.Errorf("expected %v, got %v", testStruct1, *r)
	}
}

func TestStructPtrOutput2(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct2", reflect.TypeFor[TestStruct2]())

	result, err := f.InvokeFunction("testStructPtrOutput2")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct2); !ok {
		t.Errorf("expected a pointer to a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct2, *r) {
		t.Errorf("expected %v, got %v", testStruct2, *r)
	}
}

func TestStructOutput1_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructOutput1")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testStruct1AsMap, r) {
		t.Errorf("expected %v, got %v", testStruct1AsMap, r)
	}
}

func TestStructOutput2_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructOutput2")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testStruct2AsMap, r) {
		t.Errorf("expected %v, got %v", testStruct2AsMap, r)
	}
}

func TestStructPtrOutput1_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput1")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testStruct1AsMap, r) {
		t.Errorf("expected %v, got %v", testStruct1AsMap, r)
	}
}

func TestStructPtrOutput2_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput2")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testStruct2AsMap, r) {
		t.Errorf("expected %v, got %v", testStruct2AsMap, r)
	}
}

func TestStructPtrOutput1_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput1_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestStructPtrOutput2_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput2_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}
