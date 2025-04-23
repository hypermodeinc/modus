/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript_test

import (
	"fmt"
	"maps"
	"math"
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

func TestMapInput_u8_string(t *testing.T) {
	fnName := "testMapInput_u8_string"
	m := map[uint8]string{
		1: "a",
		2: "b",
		3: "c",
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

func TestMapOutput_u8_string(t *testing.T) {
	fnName := "testMapOutput_u8_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[uint8]string{
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

func TestNullableMapInput_string_string(t *testing.T) {
	fnName := "testNullableMapInput_string_string"
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

func TestNullableMapOutput_string_string(t *testing.T) {
	fnName := "testNullableMapOutput_string_string"
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

type TestClassWithMap1 struct {
	M map[string]string
}

type TestClassWithMap2 struct {
	M map[string]any
}

func TestClassContainingMapInput_string_string(t *testing.T) {
	fnName := "testClassContainingMapInput_string_string"
	s1 := TestClassWithMap1{M: map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}}
	if _, err := fixture.CallFunction(t, fnName, s1); err != nil {
		t.Error(err)
	}

	s2 := TestClassWithMap2{M: map[string]any{
		"a": any("1"),
		"b": any("2"),
		"c": any("3"),
	}}
	if _, err := fixture.CallFunction(t, fnName, s2); err != nil {
		t.Error(err)
	}
}

func TestClassContainingMapOutput_string_string(t *testing.T) {
	fnName := "testClassContainingMapOutput_string_string"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := TestClassWithMap1{M: map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestClassWithMap1); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func makeTestMap(size int) map[string]string {
	m := make(map[string]string, size)
	for i := range size {
		key := fmt.Sprintf("key_%03d", i)
		val := fmt.Sprintf("val_%03d", i)
		m[key] = val
	}
	return m
}

var epsilon = math.Nextafter(1, 2) - 1

func TestMapInput_string_f64(t *testing.T) {
	fnName := "testMapInput_string_f64"
	m := map[string]float64{
		"a": 0.5,
		"b": epsilon,
		"c": math.SmallestNonzeroFloat64,
		"d": math.MaxFloat64,
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

func TestMapOutput_string_f64(t *testing.T) {
	fnName := "testMapOutput_string_f64"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]float64{
		"a": 0.5,
		"b": epsilon,
		"c": math.SmallestNonzeroFloat64,
		"d": math.MaxFloat64,
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]float64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !maps.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}
