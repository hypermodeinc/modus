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
	"slices"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

func TestInt8ArrayInput(t *testing.T) {
	fnName := "testInt8ArrayInput"
	arr := []int8{0, 1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestInt8ArrayBufferOutput(t *testing.T) {
	fnName := "testInt8ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []int8{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestInt16ArrayInput(t *testing.T) {
	fnName := "testInt16ArrayInput"
	arr := []int16{0, 1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestInt16ArrayBufferOutput(t *testing.T) {
	fnName := "testInt16ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []int16{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int16); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestInt32ArrayInput(t *testing.T) {
	fnName := "testInt32ArrayInput"
	arr := []int32{0, 1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestInt32ArrayBufferOutput(t *testing.T) {
	fnName := "testInt32ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []int32{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestInt64ArrayInput(t *testing.T) {
	fnName := "testInt64ArrayInput"
	arr := []int64{0, 1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestInt64ArrayBufferOutput(t *testing.T) {
	fnName := "testInt64ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []int64{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]int64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUint8ArrayInput(t *testing.T) {
	fnName := "testUint8ArrayInput"
	arr := []uint8{0, 1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestUint8ArrayInput_empty(t *testing.T) {
	fnName := "testUint8ArrayInput_empty"
	arr := []uint8{}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestUint8ArrayInput_null(t *testing.T) {
	fnName := "testUint8ArrayInput_null"
	var arr []int8 = nil

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestUint8ArrayBufferOutput(t *testing.T) {
	fnName := "testUint8ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []uint8{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]uint8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUint8ArrayBufferOutput_empty(t *testing.T) {
	fnName := "testUint8ArrayOutput_empty"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []uint8{}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]uint8); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUint8ArrayBufferOutput_null(t *testing.T) {
	fnName := "testUint8ArrayOutput_null"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestUint16ArrayInput(t *testing.T) {
	fnName := "testUint16ArrayInput"
	arr := []uint16{0, 1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestUint16ArrayBufferOutput(t *testing.T) {
	fnName := "testUint16ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []uint16{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]uint16); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUint32ArrayInput(t *testing.T) {
	fnName := "testUint32ArrayInput"
	arr := []uint32{0, 1, 2, 3}

	_, err := fixture.CallFunction(t, fnName, arr)
	if err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestUint32ArrayBufferOutput(t *testing.T) {
	fnName := "testUint32ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []uint32{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]uint32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUint64ArrayInput(t *testing.T) {
	fnName := "testUint64ArrayInput"
	arr := []uint64{0, 1, 2, 3}

	_, err := fixture.CallFunction(t, fnName, arr)
	if err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestUint64ArrayBufferOutput(t *testing.T) {
	fnName := "testUint64ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []uint64{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]uint64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestFloat32ArrayInput(t *testing.T) {
	fnName := "testFloat32ArrayInput"
	arr := []float32{0, 1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestFloat32ArrayBufferOutput(t *testing.T) {
	fnName := "testFloat32ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []float32{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]float32); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestFloat64ArrayInput(t *testing.T) {
	fnName := "testFloat64ArrayInput"
	arr := []float64{0, 1, 2, 3}

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestFloat64ArrayBufferOutput(t *testing.T) {
	fnName := "testFloat64ArrayOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []float64{0, 1, 2, 3}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]float64); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}
