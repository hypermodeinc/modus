/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"math"
	"testing"

	"hypruntime/utils"
)

func TestBoolInput_false(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testBoolInput_false", false); err != nil {
		t.Fatal(err)
	}
}

func TestBoolOutput_false(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testBoolOutput_false")
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
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testBoolInput_true", true); err != nil {
		t.Fatal(err)
	}
}

func TestBoolOutput_true(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testBoolOutput_true")
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

func TestBoolPtrInput_false(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := false
	if _, err := f.CallFunction("testBoolPtrInput_false", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testBoolPtrInput_false", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBoolPtrOutput_false(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testBoolPtrOutput_false")
	if err != nil {
		t.Fatal(err)
	}

	expected := false
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*bool); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestBoolPtrInput_true(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := true
	if _, err := f.CallFunction("testBoolPtrInput_true", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testBoolPtrInput_true", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBoolPtrOutput_true(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testBoolPtrOutput_true")
	if err != nil {
		t.Fatal(err)
	}

	expected := true
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*bool); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestBoolPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testBoolPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestBoolPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testBoolPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestByteInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testByteInput_min", byte(0)); err != nil {
		t.Fatal(err)
	}
}

func TestByteOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testByteOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := byte(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(byte); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestByteInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testByteInput_max", byte(math.MaxUint8)); err != nil {
		t.Fatal(err)
	}
}

func TestByteOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testByteOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := byte(math.MaxUint8)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(byte); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestBytePtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := byte(0)
	if _, err := f.CallFunction("testBytePtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testBytePtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBytePtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testBytePtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := byte(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*byte); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestBytePtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := byte(math.MaxUint8)
	if _, err := f.CallFunction("testBytePtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testBytePtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBytePtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testBytePtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := byte(math.MaxUint8)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*byte); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestBytePtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testBytePtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestBytePtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testBytePtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestRuneInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testRuneInput_min", rune(math.MinInt16)); err != nil {
		t.Fatal(err)
	}
}

func TestRuneOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testRuneOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := rune(math.MinInt16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(rune); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestRuneInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testRuneInput_max", rune(math.MaxInt16)); err != nil {
		t.Fatal(err)
	}
}

func TestRuneOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testRuneOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := rune(math.MaxInt16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(rune); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestRunePtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := rune(math.MinInt16)
	if _, err := f.CallFunction("testRunePtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testRunePtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestRunePtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testRunePtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := rune(math.MinInt16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*rune); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestRunePtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := rune(math.MaxInt16)
	if _, err := f.CallFunction("testRunePtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testRunePtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestRunePtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testRunePtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := rune(math.MaxInt16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*rune); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestRunePtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testRunePtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestRunePtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testRunePtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestIntInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testIntInput_min", int(math.MinInt32)); err != nil {
		t.Fatal(err)
	}
}

func TestIntOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testIntOutput_min")
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

func TestIntInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testIntInput_max", int(math.MaxInt32)); err != nil {
		t.Fatal(err)
	}
}

func TestIntOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testIntOutput_max")
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

func TestIntPtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int(math.MinInt32)
	if _, err := f.CallFunction("testIntPtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testIntPtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestIntPtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testIntPtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := int(math.MinInt32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestIntPtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int(math.MaxInt32)
	if _, err := f.CallFunction("testIntPtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testIntPtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestIntPtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testIntPtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := int(math.MaxInt32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestIntPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testIntPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testIntPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt8Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt8Input_min", int8(math.MinInt8)); err != nil {
		t.Fatal(err)
	}
}

func TestInt8Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt8Output_min")
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

func TestInt8Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt8Input_max", int8(math.MaxInt8)); err != nil {
		t.Fatal(err)
	}
}

func TestInt8Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt8Output_max")
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

func TestInt8PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int8(math.MinInt8)
	if _, err := f.CallFunction("testInt8PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testInt8PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt8PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt8PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := int8(math.MinInt8)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int8); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestInt8PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int8(math.MaxInt8)
	if _, err := f.CallFunction("testInt8PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testInt8PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt8PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt8PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := int8(math.MaxInt8)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int8); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestInt8PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt8PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt8PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt8PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt16Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt16Input_min", int16(math.MinInt16)); err != nil {
		t.Fatal(err)
	}
}

func TestInt16Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt16Output_min")
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

func TestInt16Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt16Input_max", int16(math.MaxInt16)); err != nil {
		t.Fatal(err)
	}
}

func TestInt16Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt16Output_max")
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

func TestInt16PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int16(math.MinInt16)
	if _, err := f.CallFunction("testInt16PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testInt16PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt16PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt16PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := int16(math.MinInt16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int16); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestInt16PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int16(math.MaxInt16)
	if _, err := f.CallFunction("testInt16PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testInt16PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt16PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt16PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := int16(math.MaxInt16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int16); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestInt16PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt16PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt16PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt16PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt32Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt32Input_min", int32(math.MinInt32)); err != nil {
		t.Fatal(err)
	}
}

func TestInt32Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt32Output_min")
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

func TestInt32Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt32Input_max", int32(math.MaxInt32)); err != nil {
		t.Fatal(err)
	}
}

func TestInt32Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt32Output_max")
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

func TestInt32PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int32(math.MinInt32)
	if _, err := f.CallFunction("testInt32PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testInt32PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt32PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt32PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := int32(math.MinInt32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int32); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestInt32PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int32(math.MaxInt32)
	if _, err := f.CallFunction("testInt32PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testInt32PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt32PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt32PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := int32(math.MaxInt32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int32); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestInt32PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt32PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt32PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt32PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt64Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt64Input_min", int64(math.MinInt64)); err != nil {
		t.Fatal(err)
	}
}

func TestInt64Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt64Output_min")
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

func TestInt64Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt64Input_max", int64(math.MaxInt64)); err != nil {
		t.Fatal(err)
	}
}

func TestInt64Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt64Output_max")
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

func TestInt64PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int64(math.MinInt64)
	if _, err := f.CallFunction("testInt64PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testInt64PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt64PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt64PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := int64(math.MinInt64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int64); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestInt64PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := int64(math.MaxInt64)
	if _, err := f.CallFunction("testInt64PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testInt64PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt64PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt64PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := int64(math.MaxInt64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int64); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestInt64PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testInt64PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt64PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt64PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUintInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUintInput_min", uint(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUintOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintOutput_min")
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

func TestUintInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUintInput_max", uint(math.MaxUint32)); err != nil {
		t.Fatal(err)
	}
}

func TestUintOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintOutput_max")
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

func TestUintPtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint(0)
	if _, err := f.CallFunction("testUintPtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUintPtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUintPtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintPtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUintPtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint(math.MaxUint32)
	if _, err := f.CallFunction("testUintPtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUintPtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUintPtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintPtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint(math.MaxUint32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUintPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUintPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUintPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint8Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint8Input_min", uint8(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUint8Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint8Output_min")
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

func TestUint8Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint8Input_max", uint8(math.MaxUint8)); err != nil {
		t.Fatal(err)
	}
}

func TestUint8Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint8Output_max")
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

func TestUint8PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint8(0)
	if _, err := f.CallFunction("testUint8PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUint8PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint8PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint8PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint8(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint8); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUint8PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint8(math.MaxUint8)
	if _, err := f.CallFunction("testUint8PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUint8PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint8PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint8PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint8(math.MaxUint8)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint8); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUint8PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint8PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint8PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint8PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint16Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint16Input_min", uint16(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUint16Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint16Output_min")
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

func TestUint16Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint16Input_max", uint16(math.MaxUint16)); err != nil {
		t.Fatal(err)
	}
}

func TestUint16Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint16Output_max")
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

func TestUint16PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint16(0)
	if _, err := f.CallFunction("testUint16PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUint16PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint16PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint16PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint16(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint16); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUint16PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint16(math.MaxUint16)
	if _, err := f.CallFunction("testUint16PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUint16PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint16PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint16PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint16(math.MaxUint16)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint16); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUint16PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint16PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint16PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint16PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint32Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint32Input_min", uint32(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUint32Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint32Output_min")
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

func TestUint32Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint32Input_max", uint32(math.MaxUint32)); err != nil {
		t.Fatal(err)
	}
}

func TestUint32Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint32Output_max")
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

func TestUint32PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint32(0)
	if _, err := f.CallFunction("testUint32PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUint32PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint32PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint32PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint32(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint32); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUint32PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint32(math.MaxUint32)
	if _, err := f.CallFunction("testUint32PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUint32PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint32PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint32PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint32(math.MaxUint32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint32); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUint32PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint32PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint32PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint32PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint64Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint64Input_min", uint64(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUint64Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint64Output_min")
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

func TestUint64Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint64Input_max", uint64(math.MaxUint64)); err != nil {
		t.Fatal(err)
	}
}

func TestUint64Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint64Output_max")
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

func TestUint64PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint64(0)
	if _, err := f.CallFunction("testUint64PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUint64PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint64PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint64PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint64(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint64); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUint64PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uint64(math.MaxUint64)
	if _, err := f.CallFunction("testUint64PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUint64PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint64PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint64PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := uint64(math.MaxUint64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint64); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUint64PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUint64PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint64PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint64PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUintptrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUintptrInput_min", uintptr(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintptrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := uintptr(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uintptr); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUintptrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUintptrInput_max", uintptr(math.MaxUint32)); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintptrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := uintptr(math.MaxUint32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uintptr); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestUintptrPtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uintptr(0)
	if _, err := f.CallFunction("testUintptrPtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUintptrPtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrPtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintptrPtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := uintptr(0)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uintptr); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUintptrPtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := uintptr(math.MaxUint32)
	if _, err := f.CallFunction("testUintptrPtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testUintptrPtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrPtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintptrPtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := uintptr(math.MaxUint32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uintptr); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestUintptrPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testUintptrPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUintptrPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestFloat32Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testFloat32Input_min", float32(math.SmallestNonzeroFloat32)); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat32Output_min")
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

func TestFloat32Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testFloat32Input_max", float32(math.MaxFloat32)); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat32Output_max")
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

func TestFloat32PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := float32(math.SmallestNonzeroFloat32)
	if _, err := f.CallFunction("testFloat32PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testFloat32PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat32PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := float32(math.SmallestNonzeroFloat32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*float32); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestFloat32PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := float32(math.MaxFloat32)
	if _, err := f.CallFunction("testFloat32PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testFloat32PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat32PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := float32(math.MaxFloat32)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*float32); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestFloat32PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testFloat32PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat32PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestFloat64Input_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testFloat64Input_min", float64(math.SmallestNonzeroFloat64)); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64Output_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat64Output_min")
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

func TestFloat64Input_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testFloat64Input_max", float64(math.MaxFloat64)); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64Output_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat64Output_max")
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

func TestFloat64PtrInput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := float64(math.SmallestNonzeroFloat64)
	if _, err := f.CallFunction("testFloat64PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testFloat64PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64PtrOutput_min(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat64PtrOutput_min")
	if err != nil {
		t.Fatal(err)
	}

	expected := float64(math.SmallestNonzeroFloat64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*float64); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestFloat64PtrInput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := float64(math.MaxFloat64)
	if _, err := f.CallFunction("testFloat64PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testFloat64PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64PtrOutput_max(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat64PtrOutput_max")
	if err != nil {
		t.Fatal(err)
	}

	expected := float64(math.MaxFloat64)
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*float64); !ok {
		t.Errorf("expected *%T, got %T", expected, result)
	} else if *r != expected {
		t.Errorf("expected %v, got %v", expected, *r)
	}
}

func TestFloat64PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testFloat64PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat64PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
