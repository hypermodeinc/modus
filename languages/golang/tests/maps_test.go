/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"fmt"
	"maps"
	"reflect"
	"testing"
)

func TestMapInput_string_string(t *testing.T) {
	var val = map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	if _, err := fixture.CallFunction(t, "testMapInput_string_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestMapPtrInput_string_string(t *testing.T) {
	var val = map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	if _, err := fixture.CallFunction(t, "testMapPtrInput_string_string", val); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testMapPtrInput_string_string", &val); err != nil {
		t.Fatal(err)
	}
}

func TestMapOutput_string_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testMapOutput_string_string")
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

func TestMapPtrOutput_string_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testMapPtrOutput_string_string")
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
	} else if r, ok := result.(*map[string]string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !maps.Equal(expected, *r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestIterateMap_string_string(t *testing.T) {
	var m = makeTestMap(100)

	if _, err := fixture.CallFunction(t, "testIterateMap_string_string", m); err != nil {
		t.Fatal(err)
	}
}

func TestMapLookup_string_string(t *testing.T) {
	var m = makeTestMap(100)

	result, err := fixture.CallFunction(t, "testMapLookup_string_string", m, "key_047")
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

type TestStructWithMap struct {
	M map[string]string
}

func TestStructContainingMapInput_string_string(t *testing.T) {
	s := TestStructWithMap{M: map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}}

	if _, err := fixture.CallFunction(t, "testStructContainingMapInput_string_string", s); err != nil {
		t.Fatal(err)
	}
}

func TestStructContainingMapOutput_string_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStructContainingMapOutput_string_string")
	if err != nil {
		t.Fatal(err)
	}

	expected := TestStructWithMap{M: map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStructWithMap); !ok {
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

func TestMapInput_int_float32(t *testing.T) {
	var val = map[int]float32{
		1: 1.1,
		2: 2.2,
		3: 3.3,
	}

	if _, err := fixture.CallFunction(t, "testMapInput_int_float32", val); err != nil {
		t.Fatal(err)
	}
}

func TestMapOutput_int_float32(t *testing.T) {
	result, err := fixture.CallFunction(t, "testMapOutput_int_float32")
	if err != nil {
		t.Fatal(err)
	}

	var expected = map[int]float32{
		1: 1.1,
		2: 2.2,
		3: 3.3,
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[int]float32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !maps.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}
