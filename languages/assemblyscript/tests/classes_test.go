/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"maps"
	"reflect"
	"testing"
)

type TestClass1 struct {
	A bool
}

type TestClass2 struct {
	A bool
	B int
}

type TestClass3 struct {
	A bool
	B int
	C string
}

type TestClass4 struct {
	A bool
	B int
	C *string
}

type TestClass5 struct {
	A bool
	B *TestClass3
}

var testClass1 = TestClass1{
	A: true,
}

var testClass2 = TestClass2{
	A: true,
	B: 123,
}

var testClass3 = TestClass3{
	A: true,
	B: 123,
	C: "abc",
}

var testClass4 = TestClass4{
	A: true,
	B: 123,
	C: func() *string { s := "abc"; return &s }(),
}

var testClass4_withNull = TestClass4{
	A: true,
	B: 123,
	C: nil,
}

var testClass5 = TestClass5{
	A: true,
	B: &testClass3,
}

var testClass1AsMap = map[string]any{
	"a": true,
}

var testClass2AsMap = map[string]any{
	"a": true,
	"b": 123,
}

var testClass3AsMap = map[string]any{
	"a": true,
	"b": 123,
	"c": "abc",
}

var testClass4AsMap = map[string]any{
	"a": true,
	"b": 123,
	"c": "abc",
}

var testClass4AsMap_withNull = map[string]any{
	"a": true,
	"b": 123,
	"c": nil,
}

var testClass5AsMap = map[string]any{
	"a": true,
	"b": testClass3AsMap,
}

func TestClassInput1(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput1", testClass1); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput2(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput2", testClass2); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput3(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput3", testClass3); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput4(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput4", testClass4); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput4_withNull(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput4_withNull", testClass4_withNull); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput5(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput5", testClass5); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput1_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput1", testClass1AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput2_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput2", testClass2AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput3_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput3", testClass3AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput4_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput4", testClass4AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput4_map_withNull(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput4_withNull", testClass4AsMap_withNull); err != nil {
		t.Fatal(err)
	}
}

func TestClassInput5_map(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testClassInput5", testClass5AsMap); err != nil {
		t.Fatal(err)
	}
}

func TestClassOutput1(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput1")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestClass1); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testClass1, r) {
		t.Errorf("expected %v, got %v", testClass1, r)
	}
}

func TestClassOutput2(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput2")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestClass2); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testClass2, r) {
		t.Errorf("expected %v, got %v", testClass2, r)
	}
}

func TestClassOutput3(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput3")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestClass3); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testClass3, r) {
		t.Errorf("expected %v, got %v", testClass3, r)
	}
}

func TestClassOutput4(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput4")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestClass4); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testClass4, r) {
		t.Errorf("expected %v, got %v", testClass4, r)
	}
}

func TestClassOutput4_withNull(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput4_withNull")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestClass4); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testClass4_withNull, r) {
		t.Errorf("expected %v, got %v", testClass4_withNull, r)
	}
}

func TestClassOutput5(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput5")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestClass5); !ok {
		t.Errorf("expected a struct, got %T", result)
	} else if !reflect.DeepEqual(testClass5, r) {
		t.Errorf("expected %v, got %v", testClass5, r)
	}
}

func TestClassOutput1_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput1_map")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testClass1AsMap, r) {
		t.Errorf("expected %v, got %v", testClass1AsMap, r)
	}
}

func TestClassOutput2_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput2_map")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testClass2AsMap, r) {
		t.Errorf("expected %v, got %v", testClass2AsMap, r)
	}
}

func TestClassOutput3_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput3_map")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !maps.Equal(testClass3AsMap, r) {
		t.Errorf("expected %v, got %v", testClass3AsMap, r)
	}
}

func TestClassOutput4_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput4_map")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !reflect.DeepEqual(testClass4AsMap, r) {
		t.Errorf("expected %v, got %v", testClass4AsMap, r)
	}
}

func TestClassOutput4_map_withNull(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput4_map_withNull")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !reflect.DeepEqual(testClass4AsMap_withNull, r) {
		t.Errorf("expected %v, got %v", testClass4AsMap_withNull, r)
	}
}

func TestClassOutput5_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testClassOutput5_map")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if !reflect.DeepEqual(testClass5AsMap, r) {
		t.Errorf("expected %v, got %v", testClass5AsMap, r)
	}
}
