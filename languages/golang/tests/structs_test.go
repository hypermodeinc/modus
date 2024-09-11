/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"reflect"
	"testing"

	"hypruntime/utils"
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
	"a": true,
}

var testStruct2AsMap = map[string]any{
	"a": true,
	"b": 123,
}

var testStruct3AsMap = map[string]any{
	"a": true,
	"b": 123,
	"c": "abc",
}

var testStruct4AsMap = map[string]any{
	"a": true,
	"b": 123,
	"c": func() *string { s := "abc"; return &s }(),
}

var testStruct4AsMap_withNil = map[string]any{
	"a": true,
	"b": 123,
	"c": nil,
}

func TestStructInput1(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput1", testStruct1); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput2(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput2", testStruct2); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput3(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput3", testStruct3); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput4(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput4", testStruct4); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput4_withNil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput4_withNil", testStruct4_withNil); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput1(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput1", testStruct1); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput1", &testStruct1); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput2(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput2", testStruct2); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput2", &testStruct2); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput3(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput3", testStruct3); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput3", &testStruct3); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput4", testStruct4); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput4", &testStruct4); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4_withNil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput4_withNil", testStruct4_withNil); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput4_withNil", &testStruct4_withNil); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput1_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput1", testStruct1AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput2_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput2", testStruct2AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput3_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput3", testStruct3AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput4_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput4", testStruct4AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructInput4_map_withNil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructInput4_withNil", testStruct4AsMap_withNil); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput1_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput1", testStruct1AsMap); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput1", &testStruct1AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput2_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput2", testStruct2AsMap); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput2", &testStruct2AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput3_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput3", testStruct3AsMap); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput3", &testStruct3AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput4", testStruct4AsMap); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput4", &testStruct4AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4_map_withNil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput4_withNil", testStruct4AsMap_withNil); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStructPtrInput4_withNil", &testStruct4AsMap_withNil); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput1_nil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput1_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput2_nil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput2_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput3_nil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput3_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestStructPtrInput4_nil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStructPtrInput4_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestStructOutput1(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput1")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct1

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct1); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput2(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput2")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct2

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct2); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput3(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput3")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct3

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct3); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput4(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput4")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct4

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct4); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput4_withNil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput4_withNil")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct4_withNil

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStruct4); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput1(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput1")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct1

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct1); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput2(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput2")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct2

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct2); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput3(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput3")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct3

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct3); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput4(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput4")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct4

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct4); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput4_withNil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput4_withNil")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct4_withNil

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*TestStruct4); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput1_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput1_map")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct1AsMap

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput2_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput2_map")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct2AsMap

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput3_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput3_map")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct3AsMap

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput4_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput4_map")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct4AsMap

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructOutput4_map_withNil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructOutput4_map_withNil")
	if err != nil {
		t.Fatal(err)
	}

	expected := testStruct4AsMap_withNil

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput1_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput1_map")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct1AsMap

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput2_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput2_map")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct2AsMap

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput3_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput3_map")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct3AsMap

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput4_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput4_map")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct4AsMap

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput4_map_withNil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput4_map_withNil")
	if err != nil {
		t.Fatal(err)
	}

	expected := &testStruct4AsMap_withNil

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*map[string]any); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestStructPtrOutput1_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput1_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestStructPtrOutput2_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput2_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestStructPtrOutput3_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput3_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestStructPtrOutput4_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructPtrOutput4_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
