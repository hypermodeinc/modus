/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript_test

import (
	"math"
	"testing"
)

func TestBoolInput_false(t *testing.T) {
	fnName := "testBoolInput_false"
	if _, err := fixture.CallFunction(t, fnName, false); err != nil {
		t.Error(err)
	}
}

func TestBoolOutput_false(t *testing.T) {
	fnName := "testBoolOutput_false"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := false
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(bool); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestBoolInput_true(t *testing.T) {
	fnName := "testBoolInput_true"
	if _, err := fixture.CallFunction(t, fnName, true); err != nil {
		t.Error(err)
	}
}

func TestBoolOutput_true(t *testing.T) {
	fnName := "testBoolOutput_true"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := true
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(bool); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestI8Input_min(t *testing.T) {
	fnName := "testI8Input_min"
	if _, err := fixture.CallFunction(t, fnName, int8(math.MinInt8)); err != nil {
		t.Error(err)
	}
}

func TestI8Output_min(t *testing.T) {
	fnName := "testI8Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int8(math.MinInt8)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestI8Input_max(t *testing.T) {
	fnName := "testI8Input_max"
	if _, err := fixture.CallFunction(t, fnName, int8(math.MaxInt8)); err != nil {
		t.Error(err)
	}
}

func TestI8Output_max(t *testing.T) {
	fnName := "testI8Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int8(math.MaxInt8)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestI16Input_min(t *testing.T) {
	fnName := "testI16Input_min"
	if _, err := fixture.CallFunction(t, fnName, int16(math.MinInt16)); err != nil {
		t.Error(err)
	}
}

func TestI16Output_min(t *testing.T) {
	fnName := "testI16Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int16(math.MinInt16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int16); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestI16Input_max(t *testing.T) {
	fnName := "testI16Input_max"
	if _, err := fixture.CallFunction(t, fnName, int16(math.MaxInt16)); err != nil {
		t.Error(err)
	}
}

func TestI16Output_max(t *testing.T) {
	fnName := "testI16Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int16(math.MaxInt16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int16); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestI32Input_min(t *testing.T) {
	fnName := "testI32Input_min"
	if _, err := fixture.CallFunction(t, fnName, int32(math.MinInt32)); err != nil {
		t.Error(err)
	}
}

func TestI32Output_min(t *testing.T) {
	fnName := "testI32Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int32(math.MinInt32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestI32Input_max(t *testing.T) {
	fnName := "testI32Input_max"
	if _, err := fixture.CallFunction(t, fnName, int32(math.MaxInt32)); err != nil {
		t.Error(err)
	}
}

func TestI32Output_max(t *testing.T) {
	fnName := "testI32Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int32(math.MaxInt32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestI64Input_min(t *testing.T) {
	fnName := "testI64Input_min"
	if _, err := fixture.CallFunction(t, fnName, int64(math.MinInt64)); err != nil {
		t.Error(err)
	}
}

func TestI64Output_min(t *testing.T) {
	fnName := "testI64Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int64(math.MinInt64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestI64Input_max(t *testing.T) {
	fnName := "testI64Input_max"
	if _, err := fixture.CallFunction(t, fnName, int64(math.MaxInt64)); err != nil {
		t.Error(err)
	}
}

func TestI64Output_max(t *testing.T) {
	fnName := "testI64Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int64(math.MaxInt64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestISizeInput_min(t *testing.T) {
	fnName := "testISizeInput_min"
	if _, err := fixture.CallFunction(t, fnName, int(math.MinInt32)); err != nil {
		t.Error(err)
	}
}

func TestISizeOutput_min(t *testing.T) {
	fnName := "testISizeOutput_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int(math.MinInt32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestISizeInput_max(t *testing.T) {
	fnName := "testISizeInput_max"
	if _, err := fixture.CallFunction(t, fnName, int(math.MaxInt32)); err != nil {
		t.Error(err)
	}
}

func TestISizeOutput_max(t *testing.T) {
	fnName := "testISizeOutput_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := int(math.MaxInt32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestU8Input_min(t *testing.T) {
	fnName := "testU8Input_min"
	if _, err := fixture.CallFunction(t, fnName, uint8(0)); err != nil {
		t.Error(err)
	}
}

func TestU8Output_min(t *testing.T) {
	fnName := "testU8Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint8(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestU8Input_max(t *testing.T) {
	fnName := "testU8Input_max"
	if _, err := fixture.CallFunction(t, fnName, uint8(math.MaxUint8)); err != nil {
		t.Error(err)
	}
}

func TestU8Output_max(t *testing.T) {
	fnName := "testU8Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint8(math.MaxUint8)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestU16Input_min(t *testing.T) {
	fnName := "testU16Input_min"
	if _, err := fixture.CallFunction(t, fnName, uint16(0)); err != nil {
		t.Error(err)
	}
}

func TestU16Output_min(t *testing.T) {
	fnName := "testU16Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint16(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint16); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestU16Input_max(t *testing.T) {
	fnName := "testU16Input_max"
	if _, err := fixture.CallFunction(t, fnName, uint16(math.MaxUint16)); err != nil {
		t.Error(err)
	}
}

func TestU16Output_max(t *testing.T) {
	fnName := "testU16Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint16(math.MaxUint16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint16); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestU32Input_min(t *testing.T) {
	fnName := "testU32Input_min"
	if _, err := fixture.CallFunction(t, fnName, uint32(0)); err != nil {
		t.Error(err)
	}
}

func TestU32Output_min(t *testing.T) {
	fnName := "testU32Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint32(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestU32Input_max(t *testing.T) {
	fnName := "testU32Input_max"
	if _, err := fixture.CallFunction(t, fnName, uint32(math.MaxUint32)); err != nil {
		t.Error(err)
	}
}

func TestU32Output_max(t *testing.T) {
	fnName := "testU32Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint32(math.MaxUint32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestU64Input_min(t *testing.T) {
	fnName := "testU64Input_min"
	if _, err := fixture.CallFunction(t, fnName, uint64(0)); err != nil {
		t.Error(err)
	}
}

func TestU64Output_min(t *testing.T) {
	fnName := "testU64Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint64(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestU64Input_max(t *testing.T) {
	fnName := "testU64Input_max"
	if _, err := fixture.CallFunction(t, fnName, uint64(math.MaxUint64)); err != nil {
		t.Error(err)
	}
}

func TestU64Output_max(t *testing.T) {
	fnName := "testU64Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint64(math.MaxUint64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUSizeInput_min(t *testing.T) {
	fnName := "testUSizeInput_min"
	if _, err := fixture.CallFunction(t, fnName, uint(0)); err != nil {
		t.Error(err)
	}
}

func TestUSizeOutput_min(t *testing.T) {
	fnName := "testUSizeOutput_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUSizeInput_max(t *testing.T) {
	fnName := "testUSizeInput_max"
	if _, err := fixture.CallFunction(t, fnName, uint(math.MaxUint32)); err != nil {
		t.Error(err)
	}
}

func TestUSizeOutput_max(t *testing.T) {
	fnName := "testUSizeOutput_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := uint(math.MaxUint32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestF32Input_min(t *testing.T) {
	fnName := "testF32Input_min"
	if _, err := fixture.CallFunction(t, fnName, float32(math.SmallestNonzeroFloat32)); err != nil {
		t.Error(err)
	}
}

func TestF32Output_min(t *testing.T) {
	fnName := "testF32Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := float32(math.SmallestNonzeroFloat32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(float32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestF32Input_max(t *testing.T) {
	fnName := "testF32Input_max"
	if _, err := fixture.CallFunction(t, fnName, float32(math.MaxFloat32)); err != nil {
		t.Error(err)
	}
}

func TestF32Output_max(t *testing.T) {
	fnName := "testF32Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := float32(math.MaxFloat32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(float32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestF64Input_min(t *testing.T) {
	fnName := "testF64Input_min"
	if _, err := fixture.CallFunction(t, fnName, float64(math.SmallestNonzeroFloat64)); err != nil {
		t.Error(err)
	}
}

func TestF64Output_min(t *testing.T) {
	fnName := "testF64Output_min"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := float64(math.SmallestNonzeroFloat64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(float64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestF64Input_max(t *testing.T) {
	fnName := "testF64Input_max"
	if _, err := fixture.CallFunction(t, fnName, float64(math.MaxFloat64)); err != nil {
		t.Error(err)
	}
}

func TestF64Output_max(t *testing.T) {
	fnName := "testF64Output_max"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := float64(math.MaxFloat64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(float64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}
