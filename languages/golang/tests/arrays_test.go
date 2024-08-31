/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"reflect"
	"testing"
)

func TestArrayInput0_string(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = [0]string{}

	if _, err := f.InvokeFunction("testArrayInput0_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput0_stringPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = [0]*string{}

	if _, err := f.InvokeFunction("testArrayInput0_stringPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput0_string(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput0_string")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput0_stringPtr")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = [0]*int{}

	if _, err := f.InvokeFunction("testArrayInput0_intPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput0_intPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput0_intPtr")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = [1]string{"abc"}

	if _, err := f.InvokeFunction("testArrayInput1_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput1_stringPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getStringPtrArray1()

	if _, err := f.InvokeFunction("testArrayInput1_stringPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput1_string(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput1_string")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput1_stringPtr")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getIntPtrArray1()

	if _, err := f.InvokeFunction("testArrayInput1_intPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput1_intPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput1_intPtr")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = [2]string{"abc", "def"}

	if _, err := f.InvokeFunction("testArrayInput2_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_stringPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getStringPtrArray2()

	if _, err := f.InvokeFunction("testArrayInput2_stringPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_struct(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getStructArray2()

	if _, err := f.InvokeFunction("testArrayInput2_struct", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_structPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getStructPtrArray2()

	if _, err := f.InvokeFunction("testArrayInput2_structPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_map(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getMapArray2()

	if _, err := f.InvokeFunction("testArrayInput2_map", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_mapPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getMapPtrArray2()

	if _, err := f.InvokeFunction("testArrayInput2_mapPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayInput2_intPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getIntPtrArray2()

	if _, err := f.InvokeFunction("testArrayInput2_intPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestArrayOutput2_intPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput2_intPtr")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput2_string")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput2_stringPtr")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct2", reflect.TypeFor[TestStruct2]())

	result, err := f.InvokeFunction("testArrayOutput2_struct")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("testdata.TestStruct2", reflect.TypeFor[TestStruct2]())

	result, err := f.InvokeFunction("testArrayOutput2_structPtr")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput2_map")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayOutput2_mapPtr")
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
