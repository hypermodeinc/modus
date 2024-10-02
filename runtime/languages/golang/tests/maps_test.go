/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"fmt"
	"maps"
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

func TestMapInput_string_string(t *testing.T) {
	fnName := "testMapInput_string_string"
	m := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	if _, err := fixture.CallFunction(t, fnName, m); err != nil {
		t.Error(err)
	}
	if m, err := utils.ConvertToMap(m); err != nil {
		t.Error(fmt.Errorf("failed conversion to interface map: %w", err))
	} else if _, err := fixture.CallFunction(t, fnName, m); err != nil {
		t.Error(err)
	}
}

func TestMapPtrInput_string_string(t *testing.T) {
	fnName := "testMapPtrInput_string_string"
	m := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	if _, err := fixture.CallFunction(t, fnName, m); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &m); err != nil {
		t.Error(err)
	}
	if m, err := utils.ConvertToMap(m); err != nil {
		t.Error(fmt.Errorf("failed conversion to interface map: %w", err))
	} else if _, err := fixture.CallFunction(t, fnName, m); err != nil {
		t.Error(err)
	} else if _, err := fixture.CallFunction(t, fnName, &m); err != nil {
		t.Error(err)
	}
}

func TestMapOutput_string_string(t *testing.T) {
	fnName := "testMapOutput_string_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
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
	fnName := "testMapPtrOutput_string_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
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
	fnName := "testIterateMap_string_string"
	m := makeTestMap(100)

	if _, err := fixture.CallFunction(t, fnName, m); err != nil {
		t.Error(err)
	}
}

func TestMapLookup_string_string(t *testing.T) {
	fnName := "testMapLookup_string_string"
	m := makeTestMap(100)

	result, err := fixture.CallFunction(t, fnName, m, "key_047")
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

type TestStructWithMap1 struct {
	M map[string]string
}

type TestStructWithMap2 struct {
	M map[string]any
}

func TestStructContainingMapInput_string_string(t *testing.T) {
	fnName := "testStructContainingMapInput_string_string"
	s1 := TestStructWithMap1{M: map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}}
	if _, err := fixture.CallFunction(t, fnName, s1); err != nil {
		t.Error(err)
	}

	s2 := TestStructWithMap2{M: map[string]any{
		"a": any("1"),
		"b": any("2"),
		"c": any("3"),
	}}

	if _, err := fixture.CallFunction(t, fnName, s2); err != nil {
		t.Error(err)
	}
}

func TestStructContainingMapOutput_string_string(t *testing.T) {
	fnName := "testStructContainingMapOutput_string_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := TestStructWithMap1{M: map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestStructWithMap1); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func makeTestMap(size int) map[string]string {
	m := make(map[string]string, size)
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("key_%03d", i)
		val := fmt.Sprintf("val_%03d", i)
		m[key] = val
	}
	return m
}

func TestMapInput_int_float32(t *testing.T) {
	fnName := "testMapInput_int_float32"
	m := map[int]float32{
		1: 1.1,
		2: 2.2,
		3: 3.3,
	}

	if _, err := fixture.CallFunction(t, fnName, m); err != nil {
		t.Error(err)
	}
	if m, err := utils.ConvertToMap(m); err != nil {
		t.Error(fmt.Errorf("failed conversion to interface map: %w", err))
	} else if _, err := fixture.CallFunction(t, fnName, m); err != nil {
		t.Error(err)
	}
}

func TestMapOutput_int_float32(t *testing.T) {
	fnName := "testMapOutput_int_float32"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]float32{
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
