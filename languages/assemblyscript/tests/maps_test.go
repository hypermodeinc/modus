/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"fmt"
	"maps"
	"reflect"
	"testing"
)

func TestMapInput_u8_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	var val = map[uint8]string{
		1: "a",
		2: "b",
		3: "c",
	}

	if _, err := f.CallFunction("testMapInput_u8_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestMapOutput_u8_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testMapOutput_u8_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = map[uint8]string{
		1: "a",
		2: "b",
		3: "c",
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[uint8]string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !maps.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestMapInput_string_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	var val = map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	if _, err := f.CallFunction("testMapInput_string_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestMapOutput_string_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testMapOutput_string_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !maps.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestNullableMapInput_string_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	var val = map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	if _, err := f.CallFunction("testNullableMapInput_string_string", val); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testNullableMapInput_string_string", &val); err != nil {
		t.Fatal(err)
	}
}

func TestNullableMapOutput_string_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testNullableMapOutput_string_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !maps.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestIterateMap_string_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	var m = makeTestMap(100)

	if _, err := f.CallFunction("testIterateMap_string_string", m); err != nil {
		t.Fatal(err)
	}
}

func TestMapLookup_string_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	var m = makeTestMap(100)

	result, err := f.CallFunction("testMapLookup_string_string", m, "key_047")
	if err != nil {
		t.Fatal(err)
	}

	expected := "val_047"

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if expected != r {
		t.Errorf("expected %s, got %s", expected, r)
	}
}

type testClassWithMap struct {
	M map[string]string
}

func TestClassContainingMapInput_string_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("assembly/maps/TestClassWithMap", reflect.TypeFor[testClassWithMap]())

	s := testClassWithMap{M: map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}}

	if _, err := f.CallFunction("testClassContainingMapInput_string_string", s); err != nil {
		t.Fatal(err)
	}
}

func TestClassContainingMapOutput_string_string(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("assembly/maps/TestClassWithMap", reflect.TypeFor[testClassWithMap]())

	result, err := f.CallFunction("testClassContainingMapOutput_string_string")
	if err != nil {
		t.Fatal(err)
	}

	expected := testClassWithMap{M: map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(testClassWithMap); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func makeTestMap(size int) map[string]string {
	var m = make(map[string]string, size)
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("key_%03d", i)
		val := fmt.Sprintf("val_%03d", i)
		m[key] = val
	}
	return m
}
