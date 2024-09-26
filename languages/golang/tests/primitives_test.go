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

func TestBoolPtrInput_false(t *testing.T) {
	fnName := "testBoolPtrInput_false"
	b := false

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestBoolPtrOutput_false(t *testing.T) {
	fnName := "testBoolPtrOutput_false"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testBoolPtrInput_true"
	b := true

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestBoolPtrOutput_true(t *testing.T) {
	fnName := "testBoolPtrOutput_true"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testBoolPtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestBoolPtrOutput_nil(t *testing.T) {
	fnName := "testBoolPtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Error(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestByteInput_min(t *testing.T) {
	fnName := "testByteInput_min"
	if _, err := fixture.CallFunction(t, fnName, byte(0)); err != nil {
		t.Error(err)
	}
}

func TestByteOutput_min(t *testing.T) {
	fnName := "testByteOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testByteInput_max"
	if _, err := fixture.CallFunction(t, fnName, byte(math.MaxUint8)); err != nil {
		t.Error(err)
	}
}

func TestByteOutput_max(t *testing.T) {
	fnName := "testByteOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testBytePtrInput_min"
	b := byte(0)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestBytePtrOutput_min(t *testing.T) {
	fnName := "testBytePtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testBytePtrInput_max"
	b := byte(math.MaxUint8)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestBytePtrOutput_max(t *testing.T) {
	fnName := "testBytePtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testBytePtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestBytePtrOutput_nil(t *testing.T) {
	fnName := "testBytePtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Error(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestRuneInput_min(t *testing.T) {
	fnName := "testRuneInput_min"
	if _, err := fixture.CallFunction(t, fnName, rune(math.MinInt16)); err != nil {
		t.Error(err)
	}
}

func TestRuneOutput_min(t *testing.T) {
	fnName := "testRuneOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testRuneInput_max"
	if _, err := fixture.CallFunction(t, fnName, rune(math.MaxInt16)); err != nil {
		t.Error(err)
	}
}

func TestRuneOutput_max(t *testing.T) {
	fnName := "testRuneOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testRunePtrInput_min"
	b := rune(math.MinInt16)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestRunePtrOutput_min(t *testing.T) {
	fnName := "testRunePtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testRunePtrInput_max"
	b := rune(math.MaxInt16)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestRunePtrOutput_max(t *testing.T) {
	fnName := "testRunePtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testRunePtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestRunePtrOutput_nil(t *testing.T) {
	fnName := "testRunePtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Error(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestIntInput_min(t *testing.T) {
	fnName := "testIntInput_min"
	if _, err := fixture.CallFunction(t, fnName, int(math.MinInt32)); err != nil {
		t.Error(err)
	}
}

func TestIntOutput_min(t *testing.T) {
	fnName := "testIntOutput_min"
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

func TestIntInput_max(t *testing.T) {
	fnName := "testIntInput_max"
	if _, err := fixture.CallFunction(t, fnName, int(math.MaxInt32)); err != nil {
		t.Error(err)
	}
}

func TestIntOutput_max(t *testing.T) {
	fnName := "testIntOutput_max"
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

func TestIntPtrInput_min(t *testing.T) {
	fnName := "testIntPtrInput_min"
	b := int(math.MinInt32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestIntPtrOutput_min(t *testing.T) {
	fnName := "testIntPtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testIntPtrInput_max"
	b := int(math.MaxInt32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestIntPtrOutput_max(t *testing.T) {
	fnName := "testIntPtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testIntPtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestIntPtrOutput_nil(t *testing.T) {
	fnName := "testIntPtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt8Input_min(t *testing.T) {
	fnName := "testInt8Input_min"
	if _, err := fixture.CallFunction(t, fnName, int8(math.MinInt8)); err != nil {
		t.Error(err)
	}
}

func TestInt8Output_min(t *testing.T) {
	fnName := "testInt8Output_min"
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

func TestInt8Input_max(t *testing.T) {
	fnName := "testInt8Input_max"
	if _, err := fixture.CallFunction(t, fnName, int8(math.MaxInt8)); err != nil {
		t.Error(err)
	}
}

func TestInt8Output_max(t *testing.T) {
	fnName := "testInt8Output_max"
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

func TestInt8PtrInput_min(t *testing.T) {
	fnName := "testInt8PtrInput_min"
	b := int8(math.MinInt8)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestInt8PtrOutput_min(t *testing.T) {
	fnName := "testInt8PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testInt8PtrInput_max"
	b := int8(math.MaxInt8)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestInt8PtrOutput_max(t *testing.T) {
	fnName := "testInt8PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testInt8PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestInt8PtrOutput_nil(t *testing.T) {
	fnName := "testInt8PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Error(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt16Input_min(t *testing.T) {
	fnName := "testInt16Input_min"
	if _, err := fixture.CallFunction(t, fnName, int16(math.MinInt16)); err != nil {
		t.Error(err)
	}
}

func TestInt16Output_min(t *testing.T) {
	fnName := "testInt16Output_min"
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

func TestInt16Input_max(t *testing.T) {
	fnName := "testInt16Input_max"
	if _, err := fixture.CallFunction(t, fnName, int16(math.MaxInt16)); err != nil {
		t.Error(err)
	}
}

func TestInt16Output_max(t *testing.T) {
	fnName := "testInt16Output_max"
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

func TestInt16PtrInput_min(t *testing.T) {
	fnName := "testInt16PtrInput_min"
	b := int16(math.MinInt16)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestInt16PtrOutput_min(t *testing.T) {
	fnName := "testInt16PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testInt16PtrInput_max"
	b := int16(math.MaxInt16)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestInt16PtrOutput_max(t *testing.T) {
	fnName := "testInt16PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testInt16PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestInt16PtrOutput_nil(t *testing.T) {
	fnName := "testInt16PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt32Input_min(t *testing.T) {
	fnName := "testInt32Input_min"
	if _, err := fixture.CallFunction(t, fnName, int32(math.MinInt32)); err != nil {
		t.Error(err)
	}
}

func TestInt32Output_min(t *testing.T) {
	fnName := "testInt32Output_min"
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

func TestInt32Input_max(t *testing.T) {
	fnName := "testInt32Input_max"
	if _, err := fixture.CallFunction(t, fnName, int32(math.MaxInt32)); err != nil {
		t.Error(err)
	}
}

func TestInt32Output_max(t *testing.T) {
	fnName := "testInt32Output_max"
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

func TestInt32PtrInput_min(t *testing.T) {
	fnName := "testInt32PtrInput_min"
	b := int32(math.MinInt32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestInt32PtrOutput_min(t *testing.T) {
	fnName := "testInt32PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testInt32PtrInput_max"
	b := int32(math.MaxInt32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestInt32PtrOutput_max(t *testing.T) {
	fnName := "testInt32PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testInt32PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestInt32PtrOutput_nil(t *testing.T) {
	fnName := "testInt32PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestInt64Input_min(t *testing.T) {
	fnName := "testInt64Input_min"
	if _, err := fixture.CallFunction(t, fnName, int64(math.MinInt64)); err != nil {
		t.Error(err)
	}
}

func TestInt64Output_min(t *testing.T) {
	fnName := "testInt64Output_min"
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

func TestInt64Input_max(t *testing.T) {
	fnName := "testInt64Input_max"
	if _, err := fixture.CallFunction(t, fnName, int64(math.MaxInt64)); err != nil {
		t.Error(err)
	}
}

func TestInt64Output_max(t *testing.T) {
	fnName := "testInt64Output_max"
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

func TestInt64PtrInput_min(t *testing.T) {
	fnName := "testInt64PtrInput_min"
	b := int64(math.MinInt64)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestInt64PtrOutput_min(t *testing.T) {
	fnName := "testInt64PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testInt64PtrInput_max"
	b := int64(math.MaxInt64)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestInt64PtrOutput_max(t *testing.T) {
	fnName := "testInt64PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testInt64PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestInt64PtrOutput_nil(t *testing.T) {
	fnName := "testInt64PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUintInput_min(t *testing.T) {
	fnName := "testUintInput_min"
	if _, err := fixture.CallFunction(t, fnName, uint(0)); err != nil {
		t.Error(err)
	}
}

func TestUintOutput_min(t *testing.T) {
	fnName := "testUintOutput_min"
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

func TestUintInput_max(t *testing.T) {
	fnName := "testUintInput_max"
	if _, err := fixture.CallFunction(t, fnName, uint(math.MaxUint32)); err != nil {
		t.Error(err)
	}
}

func TestUintOutput_max(t *testing.T) {
	fnName := "testUintOutput_max"
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

func TestUintPtrInput_min(t *testing.T) {
	fnName := "testUintPtrInput_min"
	b := uint(0)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUintPtrOutput_min(t *testing.T) {
	fnName := "testUintPtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUintPtrInput_max"
	b := uint(math.MaxUint32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUintPtrOutput_max(t *testing.T) {
	fnName := "testUintPtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUintPtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestUintPtrOutput_nil(t *testing.T) {
	fnName := "testUintPtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint8Input_min(t *testing.T) {
	fnName := "testUint8Input_min"
	if _, err := fixture.CallFunction(t, fnName, uint8(0)); err != nil {
		t.Error(err)
	}
}

func TestUint8Output_min(t *testing.T) {
	fnName := "testUint8Output_min"
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

func TestUint8Input_max(t *testing.T) {
	fnName := "testUint8Input_max"
	if _, err := fixture.CallFunction(t, fnName, uint8(math.MaxUint8)); err != nil {
		t.Error(err)
	}
}

func TestUint8Output_max(t *testing.T) {
	fnName := "testUint8Output_max"
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

func TestUint8PtrInput_min(t *testing.T) {
	fnName := "testUint8PtrInput_min"
	b := uint8(0)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUint8PtrOutput_min(t *testing.T) {
	fnName := "testUint8PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUint8PtrInput_max"
	b := uint8(math.MaxUint8)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUint8PtrOutput_max(t *testing.T) {
	fnName := "testUint8PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUint8PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestUint8PtrOutput_nil(t *testing.T) {
	fnName := "testUint8PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint16Input_min(t *testing.T) {
	fnName := "testUint16Input_min"
	if _, err := fixture.CallFunction(t, fnName, uint16(0)); err != nil {
		t.Error(err)
	}
}

func TestUint16Output_min(t *testing.T) {
	fnName := "testUint16Output_min"
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

func TestUint16Input_max(t *testing.T) {
	fnName := "testUint16Input_max"
	if _, err := fixture.CallFunction(t, fnName, uint16(math.MaxUint16)); err != nil {
		t.Error(err)
	}
}

func TestUint16Output_max(t *testing.T) {
	fnName := "testUint16Output_max"
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

func TestUint16PtrInput_min(t *testing.T) {
	fnName := "testUint16PtrInput_min"
	b := uint16(0)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUint16PtrOutput_min(t *testing.T) {
	fnName := "testUint16PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUint16PtrInput_max"
	b := uint16(math.MaxUint16)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUint16PtrOutput_max(t *testing.T) {
	fnName := "testUint16PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUint16PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestUint16PtrOutput_nil(t *testing.T) {
	fnName := "testUint16PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint32Input_min(t *testing.T) {
	fnName := "testUint32Input_min"
	if _, err := fixture.CallFunction(t, fnName, uint32(0)); err != nil {
		t.Error(err)
	}
}

func TestUint32Output_min(t *testing.T) {
	fnName := "testUint32Output_min"
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

func TestUint32Input_max(t *testing.T) {
	fnName := "testUint32Input_max"
	if _, err := fixture.CallFunction(t, fnName, uint32(math.MaxUint32)); err != nil {
		t.Error(err)
	}
}

func TestUint32Output_max(t *testing.T) {
	fnName := "testUint32Output_max"
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

func TestUint32PtrInput_min(t *testing.T) {
	fnName := "testUint32PtrInput_min"
	b := uint32(0)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUint32PtrOutput_min(t *testing.T) {
	fnName := "testUint32PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUint32PtrInput_max"
	b := uint32(math.MaxUint32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUint32PtrOutput_max(t *testing.T) {
	fnName := "testUint32PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUint32PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestUint32PtrOutput_nil(t *testing.T) {
	fnName := "testUint32PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUint64Input_min(t *testing.T) {
	fnName := "testUint64Input_min"
	if _, err := fixture.CallFunction(t, fnName, uint64(0)); err != nil {
		t.Error(err)
	}
}

func TestUint64Output_min(t *testing.T) {
	fnName := "testUint64Output_min"
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

func TestUint64Input_max(t *testing.T) {
	fnName := "testUint64Input_max"
	if _, err := fixture.CallFunction(t, fnName, uint64(math.MaxUint64)); err != nil {
		t.Error(err)
	}
}

func TestUint64Output_max(t *testing.T) {
	fnName := "testUint64Output_max"
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

func TestUint64PtrInput_min(t *testing.T) {
	fnName := "testUint64PtrInput_min"
	b := uint64(0)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUint64PtrOutput_min(t *testing.T) {
	fnName := "testUint64PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUint64PtrInput_max"
	b := uint64(math.MaxUint64)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUint64PtrOutput_max(t *testing.T) {
	fnName := "testUint64PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUint64PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestUint64PtrOutput_nil(t *testing.T) {
	fnName := "testUint64PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestUintptrInput_min(t *testing.T) {
	fnName := "testUintptrInput_min"
	if _, err := fixture.CallFunction(t, fnName, uintptr(0)); err != nil {
		t.Error(err)
	}
}

func TestUintptrOutput_min(t *testing.T) {
	fnName := "testUintptrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUintptrInput_max"
	if _, err := fixture.CallFunction(t, fnName, uintptr(math.MaxUint32)); err != nil {
		t.Error(err)
	}
}

func TestUintptrOutput_max(t *testing.T) {
	fnName := "testUintptrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUintptrPtrInput_min"
	b := uintptr(0)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUintptrPtrOutput_min(t *testing.T) {
	fnName := "testUintptrPtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUintptrPtrInput_max"
	b := uintptr(math.MaxUint32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestUintptrPtrOutput_max(t *testing.T) {
	fnName := "testUintptrPtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testUintptrPtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestUintptrPtrOutput_nil(t *testing.T) {
	fnName := "testUintptrPtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestFloat32Input_min(t *testing.T) {
	fnName := "testFloat32Input_min"
	if _, err := fixture.CallFunction(t, fnName, float32(math.SmallestNonzeroFloat32)); err != nil {
		t.Error(err)
	}
}

func TestFloat32Output_min(t *testing.T) {
	fnName := "testFloat32Output_min"
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

func TestFloat32Input_max(t *testing.T) {
	fnName := "testFloat32Input_max"
	if _, err := fixture.CallFunction(t, fnName, float32(math.MaxFloat32)); err != nil {
		t.Error(err)
	}
}

func TestFloat32Output_max(t *testing.T) {
	fnName := "testFloat32Output_max"
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

func TestFloat32PtrInput_min(t *testing.T) {
	fnName := "testFloat32PtrInput_min"
	b := float32(math.SmallestNonzeroFloat32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestFloat32PtrOutput_min(t *testing.T) {
	fnName := "testFloat32PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testFloat32PtrInput_max"
	b := float32(math.MaxFloat32)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestFloat32PtrOutput_max(t *testing.T) {
	fnName := "testFloat32PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testFloat32PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestFloat32PtrOutput_nil(t *testing.T) {
	fnName := "testFloat32PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestFloat64Input_min(t *testing.T) {
	fnName := "testFloat64Input_min"
	if _, err := fixture.CallFunction(t, fnName, float64(math.SmallestNonzeroFloat64)); err != nil {
		t.Error(err)
	}
}

func TestFloat64Output_min(t *testing.T) {
	fnName := "testFloat64Output_min"
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

func TestFloat64Input_max(t *testing.T) {
	fnName := "testFloat64Input_max"
	if _, err := fixture.CallFunction(t, fnName, float64(math.MaxFloat64)); err != nil {
		t.Error(err)
	}
}

func TestFloat64Output_max(t *testing.T) {
	fnName := "testFloat64Output_max"
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

func TestFloat64PtrInput_min(t *testing.T) {
	fnName := "testFloat64PtrInput_min"
	b := float64(math.SmallestNonzeroFloat64)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestFloat64PtrOutput_min(t *testing.T) {
	fnName := "testFloat64PtrOutput_min"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testFloat64PtrInput_max"
	b := float64(math.MaxFloat64)

	if _, err := fixture.CallFunction(t, fnName, b); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &b); err != nil {
		t.Error(err)
	}
}

func TestFloat64PtrOutput_max(t *testing.T) {
	fnName := "testFloat64PtrOutput_max"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testFloat64PtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestFloat64PtrOutput_nil(t *testing.T) {
	fnName := "testFloat64PtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
