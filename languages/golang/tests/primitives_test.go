/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"testing"
)

func TestBoolInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testBoolInput", true); err != nil {
		t.Fatal(err)
	}
}

func TestBoolPtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := true
	if _, err := f.InvokeFunction("testBoolPtrInput", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testBoolPtrInput", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBoolPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testBoolPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestBoolOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testBoolOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(bool); !ok {
		t.Errorf("expected a bool, got %T", result)
	} else if r != true {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestBoolPtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testBoolPtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*bool); !ok {
		t.Errorf("expected a *bool, got %T", result)
	} else if *r != true {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestBoolPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testBoolPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestByteInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testByteInput", byte(123)); err != nil {
		t.Fatal(err)
	}
}

func TestBytePtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	b := byte(123)
	if _, err := f.InvokeFunction("testBytePtrInput", b); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testBytePtrInput", &b); err != nil {
		t.Fatal(err)
	}
}

func TestBytePtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testBytePtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestByteOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testByteOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(byte); !ok {
		t.Errorf("expected a byte, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestBytePtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testBytePtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*byte); !ok {
		t.Errorf("expected a *byte, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestBytePtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testBytePtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestIntInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testIntInput", 123); err != nil {
		t.Fatal(err)
	}
}

func TestIntPtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := 123
	if _, err := f.InvokeFunction("testIntPtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testIntPtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestIntPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testIntPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testIntOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int); !ok {
		t.Errorf("expected a int, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestIntPtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testIntPtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int); !ok {
		t.Errorf("expected a *int, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestIntPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testIntPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestInt8Input(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testInt8Input", int8(123)); err != nil {
		t.Fatal(err)
	}
}

func TestInt8PtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := int8(123)
	if _, err := f.InvokeFunction("testInt8PtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testInt8PtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestInt8PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testInt8PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt8Output(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt8Output")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int8); !ok {
		t.Errorf("expected a int8, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestInt8PtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt8PtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int8); !ok {
		t.Errorf("expected a *int8, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestInt8PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt8PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestInt16Input(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testInt16Input", int16(123)); err != nil {
		t.Fatal(err)
	}
}

func TestInt16PtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := int16(123)
	if _, err := f.InvokeFunction("testInt16PtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testInt16PtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestInt16PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testInt16PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt16Output(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt16Output")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int16); !ok {
		t.Errorf("expected a int16, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestInt16PtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt16PtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int16); !ok {
		t.Errorf("expected a *int16, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestInt16PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt16PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestInt32Input(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testInt32Input", int32(123)); err != nil {
		t.Fatal(err)
	}
}

func TestInt32PtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := int32(123)
	if _, err := f.InvokeFunction("testInt32PtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testInt32PtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestInt32PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testInt32PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt32Output(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt32Output")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int32); !ok {
		t.Errorf("expected a int32, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestInt32PtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt32PtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int32); !ok {
		t.Errorf("expected a *int32, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestInt32PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt32PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestInt64Input(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testInt64Input", int64(123)); err != nil {
		t.Fatal(err)
	}
}

func TestInt64PtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := int64(123)
	if _, err := f.InvokeFunction("testInt64PtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testInt64PtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestInt64PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testInt64PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestInt64Output(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt64Output")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int64); !ok {
		t.Errorf("expected a int64, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestInt64PtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt64PtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*int64); !ok {
		t.Errorf("expected a *int64, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestInt64PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testInt64PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestRuneInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testRuneInput", rune(123)); err != nil {
		t.Fatal(err)
	}
}

func TestRunePtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := rune(123)
	if _, err := f.InvokeFunction("testRunePtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testRunePtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestRunePtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testRunePtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestRuneOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testRuneOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(rune); !ok {
		t.Errorf("expected a rune, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestRunePtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testRunePtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*rune); !ok {
		t.Errorf("expected a *rune, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestRunePtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testRunePtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestUintInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUintInput", uint(123)); err != nil {
		t.Fatal(err)
	}
}

func TestUintPtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := uint(123)
	if _, err := f.InvokeFunction("testUintPtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testUintPtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestUintPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUintPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUintOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUintOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint); !ok {
		t.Errorf("expected a uint, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestUintPtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUintPtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint); !ok {
		t.Errorf("expected a *uint, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestUintPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUintPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestUint8Input(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUint8Input", uint8(123)); err != nil {
		t.Fatal(err)
	}
}

func TestUint8PtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := uint8(123)
	if _, err := f.InvokeFunction("testUint8PtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testUint8PtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestUint8PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUint8PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint8Output(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint8Output")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint8); !ok {
		t.Errorf("expected a uint8, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestUint8PtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint8PtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint8); !ok {
		t.Errorf("expected a *uint8, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestUint8PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint8PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestUint16Input(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUint16Input", uint16(123)); err != nil {
		t.Fatal(err)
	}
}

func TestUint16PtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := uint16(123)
	if _, err := f.InvokeFunction("testUint16PtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testUint16PtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestUint16PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUint16PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint16Output(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint16Output")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint16); !ok {
		t.Errorf("expected a uint16, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestUint16PtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint16PtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint16); !ok {
		t.Errorf("expected a *uint16, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestUint16PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint16PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestUint32Input(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUint32Input", uint32(123)); err != nil {
		t.Fatal(err)
	}
}

func TestUint32PtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := uint32(123)
	if _, err := f.InvokeFunction("testUint32PtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testUint32PtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestUint32PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUint32PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint32Output(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint32Output")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint32); !ok {
		t.Errorf("expected a uint32, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestUint32PtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint32PtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint32); !ok {
		t.Errorf("expected a *uint32, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestUint32PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint32PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestUint64Input(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUint64Input", uint64(123)); err != nil {
		t.Fatal(err)
	}
}

func TestUint64PtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := uint64(123)
	if _, err := f.InvokeFunction("testUint64PtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testUint64PtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestUint64PtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUint64PtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUint64Output(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint64Output")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uint64); !ok {
		t.Errorf("expected a uint64, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestUint64PtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint64PtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uint64); !ok {
		t.Errorf("expected a *uint64, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestUint64PtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUint64PtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}

func TestUintptrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUintptrInput", uintptr(123)); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrPtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	n := uintptr(123)
	if _, err := f.InvokeFunction("testUintptrPtrInput", n); err != nil {
		t.Fatal(err)
	}
	if _, err := f.InvokeFunction("testUintptrPtrInput", &n); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.InvokeFunction("testUintptrPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestUintptrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUintptrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(uintptr); !ok {
		t.Errorf("expected a uintptr, got %T", result)
	} else if r != 123 {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestUintptrPtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUintptrPtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*uintptr); !ok {
		t.Errorf("expected a *uintptr, got %T", result)
	} else if *r != 123 {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestUintptrPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testUintptrPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Error("expected a nil result")
	}
}
