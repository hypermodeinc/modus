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

type TestStruct3 struct {
	A bool
	B int
	C string
}

type TestStruct4 struct {
	A bool
	B int
	C *string
}

var testStruct1 = TestStruct1{
	A: true,
}

var testStruct2 = TestStruct2{
	A: true,
	B: 123,
}

var testStruct3 = TestStruct3{
	A: true,
	B: 123,
	C: "abc",
}

var testStruct4 = TestStruct4{
	A: true,
	B: 123,
	C: func() *string { s := "abc"; return &s }(),
}

var testStruct4_withNil = TestStruct4{
	A: true,
	B: 123,
	C: nil,
}

var testStruct1AsMap = map[string]any{
	"A": true,
}

var testStruct2AsMap = map[string]any{
	"A": true,
	"B": 123,
}

var testStruct3AsMap = map[string]any{
	"A": true,
	"B": 123,
	"C": "abc",
}

var testStruct4AsMap = map[string]any{
	"A": true,
	"B": 123,
	"C": func() *string { s := "abc"; return &s }(),
}

var testStruct4AsMap_withNil = map[string]any{
	"A": true,
	"B": 123,
	"C": nil,
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

func TestStructInput3(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput3", testStruct3); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput4(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput4", testStruct4); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput4_withNil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput4_withNil", testStruct4_withNil); err != nil {
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

func TestStructPtrInput3(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput3", testStruct3); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput3", &testStruct3); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput4", testStruct4); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput4", &testStruct4); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4_withNil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput4_withNil", testStruct4_withNil); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput4_withNil", &testStruct4_withNil); err != nil {
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

func TestStructInput3_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput3", testStruct3AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput4_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput4", testStruct4AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput4_map_withNil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructInput4_withNil", testStruct4AsMap_withNil); err != nil {
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

func TestStructPtrInput3_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput3", testStruct3AsMap); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput3", &testStruct3AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput4", testStruct4AsMap); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput4", &testStruct4AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4_map_withNil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput4_withNil", testStruct4AsMap_withNil); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testStructPtrInput4_withNil", &testStruct4AsMap_withNil); err != nil {
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

func TestStructPtrInput3_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput3_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testStructPtrInput4_nil", nil); err != nil {
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

func TestStructOutput3(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct3", reflect.TypeFor[TestStruct3]())

	result, err := f.InvokeFunction("testStructOutput3")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct3); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct3, r) {
		t.Errorf("expected %v, got %v", testStruct3, r)
	}
}

func TestStructOutput4(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct4", reflect.TypeFor[TestStruct4]())

	result, err := f.InvokeFunction("testStructOutput4")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct4); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct4, r) {
		t.Errorf("expected %v, got %v", testStruct4, r)
	}
}

func TestStructOutput4_withNil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct4", reflect.TypeFor[TestStruct4]())

	result, err := f.InvokeFunction("testStructOutput4_withNil")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct4); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct4_withNil, r) {
		t.Errorf("expected %v, got %v", testStruct4_withNil, r)
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

func TestStructPtrOutput3(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct3", reflect.TypeFor[TestStruct3]())

	result, err := f.InvokeFunction("testStructPtrOutput3")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct3); !ok {
		t.Errorf("expected a pointer to a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct3, *r) {
		t.Errorf("expected %v, got %v", testStruct3, *r)
	}
}

func TestStructPtrOutput4(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct4", reflect.TypeFor[TestStruct4]())

	result, err := f.InvokeFunction("testStructPtrOutput4")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct4); !ok {
		t.Errorf("expected a pointer to a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct4, *r) {
		t.Errorf("expected %v, got %v", testStruct4, *r)
	}
}

func TestStructPtrOutput4_withNil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct4", reflect.TypeFor[TestStruct4]())

	result, err := f.InvokeFunction("testStructPtrOutput4_withNil")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct4); !ok {
		t.Errorf("expected a pointer to a struct, got %T", result)
	} else if !reflect.DeepEqual(testStruct4_withNil, *r) {
		t.Errorf("expected %v, got %v", testStruct4_withNil, *r)
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

func TestStructOutput3_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructOutput3")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testStruct3AsMap, r) {
		t.Errorf("expected %v, got %v", testStruct3AsMap, r)
	}
}

func TestStructOutput4_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructOutput4")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !reflect.DeepEqual(testStruct4AsMap, r) {
		t.Errorf("expected %v, got %v", testStruct4AsMap, r)
	}
}

func TestStructOutput4_map_withNil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructOutput4_withNil")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !reflect.DeepEqual(testStruct4AsMap_withNil, r) {
		t.Errorf("expected %v, got %v", testStruct4AsMap_withNil, r)
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
		t.Errorf("expected a map[string]interface{}, got %T", result)
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

func TestStructPtrOutput3_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput3")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testStruct3AsMap, r) {
		t.Errorf("expected %v, got %v", testStruct3AsMap, r)
	}
}

func TestStructPtrOutput4_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput4")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !reflect.DeepEqual(testStruct4AsMap, r) {
		t.Errorf("expected %v, got %v", testStruct4AsMap, r)
	}
}

func TestStructPtrOutput4_map_withNil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput4_withNil")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !reflect.DeepEqual(testStruct4AsMap_withNil, r) {
		t.Errorf("expected %v, got %v", testStruct4AsMap_withNil, r)
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

func TestStructPtrOutput3_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput3_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestStructPtrOutput4_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testStructPtrOutput4_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}
