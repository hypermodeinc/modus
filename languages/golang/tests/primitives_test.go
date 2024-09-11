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
	if _, err := fixture.CallFunction(t, "testBoolInput_false", false); err != nil {
		t.Fatal(err)
	}
}

func TestBoolOutput_false(t *testing.T) {
	result, err := fixture.CallFunction(t, "testBoolOutput_false")
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
	if _, err := fixture.CallFunction(t, "testBoolInput_true", true); err != nil {
		t.Fatal(err)
	}
}

func TestBoolOutput_true(t *testing.T) {
	result, err := fixture.CallFunction(t, "testBoolOutput_true")
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
	b := false
	if _, err := fixture.CallFunction(t, "testBoolPtrInput_false", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testBoolPtrInput_false", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBoolPtrOutput_false(t *testing.T) {
	result, err := fixture.CallFunction(t, "testBoolPtrOutput_false")
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
	b := true
	if _, err := fixture.CallFunction(t, "testBoolPtrInput_true", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testBoolPtrInput_true", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBoolPtrOutput_true(t *testing.T) {
	result, err := fixture.CallFunction(t, "testBoolPtrOutput_true")
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
	if _, err := fixture.CallFunction(t, "testBoolPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestBoolPtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testBoolPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestByteInput_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testByteInput_min", byte(0)); err != nil {
		t.Fatal(err)
	}
}

func TestByteOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testByteOutput_min")
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
	if _, err := fixture.CallFunction(t, "testByteInput_max", byte(math.MaxUint8)); err != nil {
		t.Fatal(err)
	}
}

func TestByteOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testByteOutput_max")
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
	b := byte(0)
	if _, err := fixture.CallFunction(t, "testBytePtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testBytePtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBytePtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testBytePtrOutput_min")
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
	b := byte(math.MaxUint8)
	if _, err := fixture.CallFunction(t, "testBytePtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testBytePtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBytePtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testBytePtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testBytePtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestBytePtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testBytePtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestRuneInput_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testRuneInput_min", rune(math.MinInt16)); err != nil {
		t.Fatal(err)
	}
}

func TestRuneOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testRuneOutput_min")
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
	if _, err := fixture.CallFunction(t, "testRuneInput_max", rune(math.MaxInt16)); err != nil {
		t.Fatal(err)
	}
}

func TestRuneOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testRuneOutput_max")
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
	b := rune(math.MinInt16)
	if _, err := fixture.CallFunction(t, "testRunePtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testRunePtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestRunePtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testRunePtrOutput_min")
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
	b := rune(math.MaxInt16)
	if _, err := fixture.CallFunction(t, "testRunePtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testRunePtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestRunePtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testRunePtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testRunePtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestRunePtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testRunePtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestIntInput_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testIntInput_min", int(math.MinInt32)); err != nil {
		t.Fatal(err)
	}
}

func TestIntOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testIntOutput_min")
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
	if _, err := fixture.CallFunction(t, "testIntInput_max", int(math.MaxInt32)); err != nil {
		t.Fatal(err)
	}
}

func TestIntOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testIntOutput_max")
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
	b := int(math.MinInt32)
	if _, err := fixture.CallFunction(t, "testIntPtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testIntPtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestIntPtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testIntPtrOutput_min")
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
	b := int(math.MaxInt32)
	if _, err := fixture.CallFunction(t, "testIntPtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testIntPtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestIntPtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testIntPtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testIntPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntPtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testIntPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt8Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testInt8Input_min", int8(math.MinInt8)); err != nil {
		t.Fatal(err)
	}
}

func TestInt8Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt8Output_min")
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
	if _, err := fixture.CallFunction(t, "testInt8Input_max", int8(math.MaxInt8)); err != nil {
		t.Fatal(err)
	}
}

func TestInt8Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt8Output_max")
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
	b := int8(math.MinInt8)
	if _, err := fixture.CallFunction(t, "testInt8PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testInt8PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt8PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt8PtrOutput_min")
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
	b := int8(math.MaxInt8)
	if _, err := fixture.CallFunction(t, "testInt8PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testInt8PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt8PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt8PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testInt8PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt8PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt8PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt16Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testInt16Input_min", int16(math.MinInt16)); err != nil {
		t.Fatal(err)
	}
}

func TestInt16Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt16Output_min")
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
	if _, err := fixture.CallFunction(t, "testInt16Input_max", int16(math.MaxInt16)); err != nil {
		t.Fatal(err)
	}
}

func TestInt16Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt16Output_max")
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
	b := int16(math.MinInt16)
	if _, err := fixture.CallFunction(t, "testInt16PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testInt16PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt16PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt16PtrOutput_min")
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
	b := int16(math.MaxInt16)
	if _, err := fixture.CallFunction(t, "testInt16PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testInt16PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt16PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt16PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testInt16PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt16PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt16PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt32Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testInt32Input_min", int32(math.MinInt32)); err != nil {
		t.Fatal(err)
	}
}

func TestInt32Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt32Output_min")
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
	if _, err := fixture.CallFunction(t, "testInt32Input_max", int32(math.MaxInt32)); err != nil {
		t.Fatal(err)
	}
}

func TestInt32Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt32Output_max")
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
	b := int32(math.MinInt32)
	if _, err := fixture.CallFunction(t, "testInt32PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testInt32PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt32PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt32PtrOutput_min")
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
	b := int32(math.MaxInt32)
	if _, err := fixture.CallFunction(t, "testInt32PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testInt32PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt32PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt32PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testInt32PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt32PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt32PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt64Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testInt64Input_min", int64(math.MinInt64)); err != nil {
		t.Fatal(err)
	}
}

func TestInt64Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt64Output_min")
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
	if _, err := fixture.CallFunction(t, "testInt64Input_max", int64(math.MaxInt64)); err != nil {
		t.Fatal(err)
	}
}

func TestInt64Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt64Output_max")
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
	b := int64(math.MinInt64)
	if _, err := fixture.CallFunction(t, "testInt64PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testInt64PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt64PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt64PtrOutput_min")
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
	b := int64(math.MaxInt64)
	if _, err := fixture.CallFunction(t, "testInt64PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testInt64PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestInt64PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt64PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testInt64PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt64PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testInt64PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUintInput_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testUintInput_min", uint(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUintOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintOutput_min")
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
	if _, err := fixture.CallFunction(t, "testUintInput_max", uint(math.MaxUint32)); err != nil {
		t.Fatal(err)
	}
}

func TestUintOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintOutput_max")
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
	b := uint(0)
	if _, err := fixture.CallFunction(t, "testUintPtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUintPtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUintPtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintPtrOutput_min")
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
	b := uint(math.MaxUint32)
	if _, err := fixture.CallFunction(t, "testUintPtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUintPtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUintPtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintPtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testUintPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUintPtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint8Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testUint8Input_min", uint8(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUint8Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint8Output_min")
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
	if _, err := fixture.CallFunction(t, "testUint8Input_max", uint8(math.MaxUint8)); err != nil {
		t.Fatal(err)
	}
}

func TestUint8Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint8Output_max")
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
	b := uint8(0)
	if _, err := fixture.CallFunction(t, "testUint8PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUint8PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint8PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint8PtrOutput_min")
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
	b := uint8(math.MaxUint8)
	if _, err := fixture.CallFunction(t, "testUint8PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUint8PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint8PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint8PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testUint8PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint8PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint8PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint16Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testUint16Input_min", uint16(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUint16Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint16Output_min")
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
	if _, err := fixture.CallFunction(t, "testUint16Input_max", uint16(math.MaxUint16)); err != nil {
		t.Fatal(err)
	}
}

func TestUint16Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint16Output_max")
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
	b := uint16(0)
	if _, err := fixture.CallFunction(t, "testUint16PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUint16PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint16PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint16PtrOutput_min")
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
	b := uint16(math.MaxUint16)
	if _, err := fixture.CallFunction(t, "testUint16PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUint16PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint16PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint16PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testUint16PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint16PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint16PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint32Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testUint32Input_min", uint32(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUint32Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint32Output_min")
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
	if _, err := fixture.CallFunction(t, "testUint32Input_max", uint32(math.MaxUint32)); err != nil {
		t.Fatal(err)
	}
}

func TestUint32Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint32Output_max")
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
	b := uint32(0)
	if _, err := fixture.CallFunction(t, "testUint32PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUint32PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint32PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint32PtrOutput_min")
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
	b := uint32(math.MaxUint32)
	if _, err := fixture.CallFunction(t, "testUint32PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUint32PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint32PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint32PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testUint32PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint32PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint32PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint64Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testUint64Input_min", uint64(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUint64Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint64Output_min")
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
	if _, err := fixture.CallFunction(t, "testUint64Input_max", uint64(math.MaxUint64)); err != nil {
		t.Fatal(err)
	}
}

func TestUint64Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint64Output_max")
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
	b := uint64(0)
	if _, err := fixture.CallFunction(t, "testUint64PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUint64PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint64PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint64PtrOutput_min")
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
	b := uint64(math.MaxUint64)
	if _, err := fixture.CallFunction(t, "testUint64PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUint64PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUint64PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint64PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testUint64PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint64PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUint64PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUintptrInput_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testUintptrInput_min", uintptr(0)); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintptrOutput_min")
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
	if _, err := fixture.CallFunction(t, "testUintptrInput_max", uintptr(math.MaxUint32)); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintptrOutput_max")
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
	b := uintptr(0)
	if _, err := fixture.CallFunction(t, "testUintptrPtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUintptrPtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrPtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintptrPtrOutput_min")
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
	b := uintptr(math.MaxUint32)
	if _, err := fixture.CallFunction(t, "testUintptrPtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testUintptrPtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrPtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintptrPtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testUintptrPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrPtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUintptrPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestFloat32Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testFloat32Input_min", float32(math.SmallestNonzeroFloat32)); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat32Output_min")
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
	if _, err := fixture.CallFunction(t, "testFloat32Input_max", float32(math.MaxFloat32)); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat32Output_max")
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
	b := float32(math.SmallestNonzeroFloat32)
	if _, err := fixture.CallFunction(t, "testFloat32PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testFloat32PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat32PtrOutput_min")
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
	b := float32(math.MaxFloat32)
	if _, err := fixture.CallFunction(t, "testFloat32PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testFloat32PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat32PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testFloat32PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestFloat32PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat32PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestFloat64Input_min(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testFloat64Input_min", float64(math.SmallestNonzeroFloat64)); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat64Output_min")
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
	if _, err := fixture.CallFunction(t, "testFloat64Input_max", float64(math.MaxFloat64)); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat64Output_max")
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
	b := float64(math.SmallestNonzeroFloat64)
	if _, err := fixture.CallFunction(t, "testFloat64PtrInput_min", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testFloat64PtrInput_min", &b); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64PtrOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat64PtrOutput_min")
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
	b := float64(math.MaxFloat64)
	if _, err := fixture.CallFunction(t, "testFloat64PtrInput_max", b); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testFloat64PtrInput_max", &b); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64PtrOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat64PtrOutput_max")
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
	if _, err := fixture.CallFunction(t, "testFloat64PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64PtrOutput_nil(t *testing.T) {
	result, err := fixture.CallFunction(t, "testFloat64PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
