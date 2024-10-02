/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

func TestArrayInput0_string(t *testing.T) {
	fnName := "testArrayInput0_string"
	arr := [0]string{}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput0_stringPtr(t *testing.T) {
	fnName := "testArrayInput0_stringPtr"
	arr := [0]*string{}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput0_string(t *testing.T) {
	fnName := "testArrayOutput0_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [0]string{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([0]string); !ok {
		t.Errorf("expected a [0]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput0_stringPtr(t *testing.T) {
	fnName := "testArrayOutput0_stringPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [0]*string{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([0]*string); !ok {
		t.Errorf("expected a [0]*string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput0_intPtr(t *testing.T) {
	fnName := "testArrayInput0_intPtr"
	arr := [0]*int{}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput0_intPtr(t *testing.T) {
	fnName := "testArrayOutput0_intPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [0]*int{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([0]*int); !ok {
		t.Errorf("expected a [0]*int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput1_string(t *testing.T) {
	fnName := "testArrayInput1_string"
	arr := [1]string{"abc"}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput1_stringPtr(t *testing.T) {
	fnName := "testArrayInput1_stringPtr"
	arr := getStringPtrArray1()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput1_string(t *testing.T) {
	fnName := "testArrayOutput1_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [1]string{"abc"}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([1]string); !ok {
		t.Errorf("expected a [1]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput1_stringPtr(t *testing.T) {
	fnName := "testArrayOutput1_stringPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getStringPtrArray1()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([1]*string); !ok {
		t.Errorf("expected a [1]*string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput1_intPtr(t *testing.T) {
	fnName := "testArrayInput1_intPtr"
	arr := getIntPtrArray1()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput1_intPtr(t *testing.T) {
	fnName := "testArrayOutput1_intPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getIntPtrArray1()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([1]*int); !ok {
		t.Errorf("expected a [1]*int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput2_string(t *testing.T) {
	fnName := "testArrayInput2_string"
	arr := [2]string{"abc", "def"}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput2_stringPtr(t *testing.T) {
	fnName := "testArrayInput2_stringPtr"
	arr := getStringPtrArray2()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput2_struct(t *testing.T) {
	fnName := "testArrayInput2_struct"
	arr := getStructArray2()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput2_structPtr(t *testing.T) {
	fnName := "testArrayInput2_structPtr"
	arr := getStructPtrArray2()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput2_map(t *testing.T) {
	fnName := "testArrayInput2_map"
	arr := getMapArray2()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput2_mapPtr(t *testing.T) {
	fnName := "testArrayInput2_mapPtr"
	arr := getMapPtrArray2()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput2_intPtr(t *testing.T) {
	fnName := "testArrayInput2_intPtr"
	arr := getIntPtrArray2()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput2_intPtr(t *testing.T) {
	fnName := "testArrayOutput2_intPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getIntPtrArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]*int); !ok {
		t.Errorf("expected a [2]*int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_string(t *testing.T) {
	fnName := "testArrayOutput2_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [2]string{"abc", "def"}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]string); !ok {
		t.Errorf("expected a [2]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_stringPtr(t *testing.T) {
	fnName := "testArrayOutput2_stringPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getStringPtrArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]*string); !ok {
		t.Errorf("expected a [2]*string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_struct(t *testing.T) {
	fnName := "testArrayOutput2_struct"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getStructArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]TestStruct2); !ok {
		t.Errorf("expected a [2]TestStruct2, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_structPtr(t *testing.T) {
	fnName := "testArrayOutput2_structPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getStructPtrArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]*TestStruct2); !ok {
		t.Errorf("expected a [2]*TestStruct2, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_map(t *testing.T) {
	fnName := "testArrayOutput2_map"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getMapArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]map[string]string); !ok {
		t.Errorf("expected a [2]map[string]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_mapPtr(t *testing.T) {
	fnName := "testArrayOutput2_mapPtr"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getMapPtrArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]*map[string]string); !ok {
		t.Errorf("expected a [2]*map[string]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestPtrArrayInput1_int(t *testing.T) {
	fnName := "testPtrArrayInput1_int"
	arr := getPtrIntArray1()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](*arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, &arr); err != nil {
		t.Error(err)
	}
}

func TestPtrArrayInput2_int(t *testing.T) {
	fnName := "testPtrArrayInput2_int"
	arr := getPtrIntArray2()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](*arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, &arr); err != nil {
		t.Error(err)
	}
}

func TestPtrArrayInput1_string(t *testing.T) {
	fnName := "testPtrArrayInput1_string"
	arr := getPtrStringArray1()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](*arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, &arr); err != nil {
		t.Error(err)
	}
}

func TestPtrArrayInput2_string(t *testing.T) {
	fnName := "testPtrArrayInput2_string"
	arr := getPtrStringArray2()

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](*arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, &arr); err != nil {
		t.Error(err)
	}
}

func TestPtrArrayOutput1_int(t *testing.T) {
	fnName := "testPtrArrayOutput1_int"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getPtrIntArray1()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*[1]int); !ok {
		t.Errorf("expected a *[1]int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestPtrArrayOutput2_int(t *testing.T) {
	fnName := "testPtrArrayOutput2_int"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getPtrIntArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*[2]int); !ok {
		t.Errorf("expected a *[2]int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestPtrArrayOutput1_string(t *testing.T) {
	fnName := "testPtrArrayOutput1_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getPtrStringArray1()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*[1]string); !ok {
		t.Errorf("expected a *[1]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestPtrArrayOutput2_string(t *testing.T) {
	fnName := "testPtrArrayOutput2_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := getPtrStringArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*[2]string); !ok {
		t.Errorf("expected a *[2]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput0_byte(t *testing.T) {
	fnName := "testArrayInput0_byte"
	arr := [0]byte{}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput1_byte(t *testing.T) {
	fnName := "testArrayInput1_byte"
	arr := [1]byte{1}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayInput2_byte(t *testing.T) {
	fnName := "testArrayInput2_byte"
	arr := [2]byte{1, 2}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayOutput0_byte(t *testing.T) {
	fnName := "testArrayOutput0_byte"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [0]byte{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([0]byte); !ok {
		t.Errorf("expected a [0]byte, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput1_byte(t *testing.T) {
	fnName := "testArrayOutput1_byte"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [1]byte{1}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([1]byte); !ok {
		t.Errorf("expected a [1]byte, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_byte(t *testing.T) {
	fnName := "testArrayOutput2_byte"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := [2]byte{1, 2}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]byte); !ok {
		t.Errorf("expected a [2]byte, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func getIntPtrArray1() [1]*int {
	a := 11
	return [1]*int{&a}
}

func getIntPtrArray2() [2]*int {
	a := 11
	b := 22
	return [2]*int{&a, &b}
}

func getStringPtrArray1() [1]*string {
	a := "abc"
	return [1]*string{&a}
}

func getStringPtrArray2() [2]*string {
	a := "abc"
	b := "def"
	return [2]*string{&a, &b}
}

func getStructArray2() [2]TestStruct2 {
	return [2]TestStruct2{
		{A: true, B: 123},
		{A: false, B: 456},
	}
}

func getStructPtrArray2() [2]*TestStruct2 {
	return [2]*TestStruct2{
		{A: true, B: 123},
		{A: false, B: 456},
	}
}

func getMapArray2() [2]map[string]string {
	return [2]map[string]string{
		{"A": "true", "B": "123"},
		{"C": "false", "D": "456"},
	}
}

func getMapPtrArray2() [2]*map[string]string {
	return [2]*map[string]string{
		{"A": "true", "B": "123"},
		{"C": "false", "D": "456"},
	}
}

func getPtrIntArray1() *[1]int {
	a := 11
	return &[1]int{a}
}

func getPtrIntArray2() *[2]int {
	a := 11
	b := 22
	return &[2]int{a, b}
}

func getPtrStringArray1() *[1]string {
	a := "abc"
	return &[1]string{a}
}

func getPtrStringArray2() *[2]string {
	a := "abc"
	b := "def"
	return &[2]string{a, b}
}
