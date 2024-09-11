/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"reflect"
	"testing"
)

func TestArrayInput0_string(t *testing.T) {
	var val = [0]string{}

	if _, err := fixture.CallFunction(t, "testArrayInput0_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput0_stringPtr(t *testing.T) {
	var val = [0]*string{}

	if _, err := fixture.CallFunction(t, "testArrayInput0_stringPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput0_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput0_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = [0]string{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([0]string); !ok {
		t.Errorf("expected a [0]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput0_stringPtr(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput0_stringPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = [0]*string{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([0]*string); !ok {
		t.Errorf("expected a [0]*string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput0_intPtr(t *testing.T) {
	var val = [0]*int{}

	if _, err := fixture.CallFunction(t, "testArrayInput0_intPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput0_intPtr(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput0_intPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = [0]*int{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([0]*int); !ok {
		t.Errorf("expected a [0]*int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput1_string(t *testing.T) {
	var val = [1]string{"abc"}

	if _, err := fixture.CallFunction(t, "testArrayInput1_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput1_stringPtr(t *testing.T) {
	var val = getStringPtrArray1()

	if _, err := fixture.CallFunction(t, "testArrayInput1_stringPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput1_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput1_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = [1]string{"abc"}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([1]string); !ok {
		t.Errorf("expected a [1]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput1_stringPtr(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput1_stringPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getStringPtrArray1()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([1]*string); !ok {
		t.Errorf("expected a [1]*string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput1_intPtr(t *testing.T) {
	var val = getIntPtrArray1()

	if _, err := fixture.CallFunction(t, "testArrayInput1_intPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput1_intPtr(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput1_intPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getIntPtrArray1()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([1]*int); !ok {
		t.Errorf("expected a [1]*int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput2_string(t *testing.T) {
	var val = [2]string{"abc", "def"}

	if _, err := fixture.CallFunction(t, "testArrayInput2_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_stringPtr(t *testing.T) {
	var val = getStringPtrArray2()

	if _, err := fixture.CallFunction(t, "testArrayInput2_stringPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_struct(t *testing.T) {
	var val = getStructArray2()

	if _, err := fixture.CallFunction(t, "testArrayInput2_struct", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_structPtr(t *testing.T) {
	var val = getStructPtrArray2()

	if _, err := fixture.CallFunction(t, "testArrayInput2_structPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_map(t *testing.T) {
	var val = getMapArray2()

	if _, err := fixture.CallFunction(t, "testArrayInput2_map", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_mapPtr(t *testing.T) {
	var val = getMapPtrArray2()

	if _, err := fixture.CallFunction(t, "testArrayInput2_mapPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_intPtr(t *testing.T) {
	var val = getIntPtrArray2()

	if _, err := fixture.CallFunction(t, "testArrayInput2_intPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput2_intPtr(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput2_intPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getIntPtrArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]*int); !ok {
		t.Errorf("expected a [2]*int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput2_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = [2]string{"abc", "def"}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]string); !ok {
		t.Errorf("expected a [2]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_stringPtr(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput2_stringPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getStringPtrArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]*string); !ok {
		t.Errorf("expected a [2]*string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_struct(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput2_struct")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getStructArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]TestStruct2); !ok {
		t.Errorf("expected a [2]TestStruct2, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_structPtr(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput2_structPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getStructPtrArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]*TestStruct2); !ok {
		t.Errorf("expected a [2]*TestStruct2, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_map(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput2_map")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getMapArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]map[string]string); !ok {
		t.Errorf("expected a [2]map[string]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_mapPtr(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput2_mapPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getMapPtrArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([2]*map[string]string); !ok {
		t.Errorf("expected a [2]*map[string]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestPtrArrayInput1_int(t *testing.T) {
	var val = getPtrIntArray1()

	if _, err := fixture.CallFunction(t, "testPtrArrayInput1_int", val); err != nil {
		t.Fatal(err)
	}
}

func TestPtrArrayInput2_int(t *testing.T) {
	var val = getPtrIntArray2()

	if _, err := fixture.CallFunction(t, "testPtrArrayInput2_int", val); err != nil {
		t.Fatal(err)
	}
}

func TestPtrArrayInput1_string(t *testing.T) {
	var val = getPtrStringArray1()

	if _, err := fixture.CallFunction(t, "testPtrArrayInput1_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestPtrArrayInput2_string(t *testing.T) {
	var val = getPtrStringArray2()

	if _, err := fixture.CallFunction(t, "testPtrArrayInput2_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestPtrArrayOutput1_int(t *testing.T) {
	result, err := fixture.CallFunction(t, "testPtrArrayOutput1_int")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getPtrIntArray1()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*[1]int); !ok {
		t.Errorf("expected a *[1]int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestPtrArrayOutput2_int(t *testing.T) {
	result, err := fixture.CallFunction(t, "testPtrArrayOutput2_int")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getPtrIntArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*[2]int); !ok {
		t.Errorf("expected a *[2]int, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestPtrArrayOutput1_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testPtrArrayOutput1_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getPtrStringArray1()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*[1]string); !ok {
		t.Errorf("expected a *[1]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestPtrArrayOutput2_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testPtrArrayOutput2_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getPtrStringArray2()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*[2]string); !ok {
		t.Errorf("expected a *[2]string, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayInput0_byte(t *testing.T) {
	var val = [0]byte{}

	if _, err := fixture.CallFunction(t, "testArrayInput0_byte", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput1_byte(t *testing.T) {
	var val = [1]byte{1}

	if _, err := fixture.CallFunction(t, "testArrayInput1_byte", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_byte(t *testing.T) {
	var val = [2]byte{1, 2}

	if _, err := fixture.CallFunction(t, "testArrayInput2_byte", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput0_byte(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput0_byte")
	if err != nil {
		t.Fatal(err)
	}

	var expected = [0]byte{}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([0]byte); !ok {
		t.Errorf("expected a [0]byte, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput1_byte(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput1_byte")
	if err != nil {
		t.Fatal(err)
	}

	var expected = [1]byte{1}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([1]byte); !ok {
		t.Errorf("expected a [1]byte, got %T", result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestArrayOutput2_byte(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayOutput2_byte")
	if err != nil {
		t.Fatal(err)
	}

	var expected = [2]byte{1, 2}
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
