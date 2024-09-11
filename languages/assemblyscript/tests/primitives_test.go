/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"math"
	"testing"
)

func TestBoolInput_false(t *testing.T) {
	_, err := fixture.CallFunction(t, "testBoolInput_false", false)
	if err != nil {
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
	_, err := fixture.CallFunction(t, "testBoolInput_true", true)
	if err != nil {
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

func TestI8Input_min(t *testing.T) {
	_, err := fixture.CallFunction(t, "testI8Input_min", int8(math.MinInt8))
	if err != nil {
		t.Fatal(err)
	}
}

func TestI8Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testI8Output_min")
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
	_, err := fixture.CallFunction(t, "testI8Input_max", int8(math.MaxInt8))
	if err != nil {
		t.Fatal(err)
	}
}

func TestI8Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testI8Output_max")
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
	_, err := fixture.CallFunction(t, "testI16Input_min", int16(math.MinInt16))
	if err != nil {
		t.Fatal(err)
	}
}

func TestI16Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testI16Output_min")
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
	_, err := fixture.CallFunction(t, "testI16Input_max", int16(math.MaxInt16))
	if err != nil {
		t.Fatal(err)
	}
}

func TestI16Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testI16Output_max")
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
	_, err := fixture.CallFunction(t, "testI32Input_min", int32(math.MinInt32))
	if err != nil {
		t.Fatal(err)
	}
}

func TestI32Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testI32Output_min")
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
	_, err := fixture.CallFunction(t, "testI32Input_max", int32(math.MaxInt32))
	if err != nil {
		t.Fatal(err)
	}
}

func TestI32Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testI32Output_max")
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
	_, err := fixture.CallFunction(t, "testI64Input_min", int64(math.MinInt64))
	if err != nil {
		t.Fatal(err)
	}
}

func TestI64Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testI64Output_min")
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
	_, err := fixture.CallFunction(t, "testI64Input_max", int64(math.MaxInt64))
	if err != nil {
		t.Fatal(err)
	}
}

func TestI64Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testI64Output_max")
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
	_, err := fixture.CallFunction(t, "testISizeInput_min", int(math.MinInt32))
	if err != nil {
		t.Fatal(err)
	}
}

func TestISizeOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testISizeOutput_min")
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
	_, err := fixture.CallFunction(t, "testISizeInput_max", int(math.MaxInt32))
	if err != nil {
		t.Fatal(err)
	}
}

func TestISizeOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testISizeOutput_max")
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
	_, err := fixture.CallFunction(t, "testU8Input_min", uint8(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestU8Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testU8Output_min")
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
	_, err := fixture.CallFunction(t, "testU8Input_max", uint8(math.MaxUint8))
	if err != nil {
		t.Fatal(err)
	}
}

func TestU8Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testU8Output_max")
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
	_, err := fixture.CallFunction(t, "testU16Input_min", uint16(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestU16Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testU16Output_min")
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
	_, err := fixture.CallFunction(t, "testU16Input_max", uint16(math.MaxUint16))
	if err != nil {
		t.Fatal(err)
	}
}

func TestU16Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testU16Output_max")
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
	_, err := fixture.CallFunction(t, "testU32Input_min", uint32(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestU32Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testU32Output_min")
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
	_, err := fixture.CallFunction(t, "testU32Input_max", uint32(math.MaxUint32))
	if err != nil {
		t.Fatal(err)
	}
}

func TestU32Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testU32Output_max")
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
	_, err := fixture.CallFunction(t, "testU64Input_min", uint64(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestU64Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testU64Output_min")
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
	_, err := fixture.CallFunction(t, "testU64Input_max", uint64(math.MaxUint64))
	if err != nil {
		t.Fatal(err)
	}
}

func TestU64Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testU64Output_max")
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
	_, err := fixture.CallFunction(t, "testUSizeInput_min", uint(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestUSizeOutput_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUSizeOutput_min")
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
	_, err := fixture.CallFunction(t, "testUSizeInput_max", uint(math.MaxUint32))
	if err != nil {
		t.Fatal(err)
	}
}

func TestUSizeOutput_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testUSizeOutput_max")
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
	_, err := fixture.CallFunction(t, "testF32Input_min", float32(math.SmallestNonzeroFloat32))
	if err != nil {
		t.Fatal(err)
	}
}

func TestF32Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testF32Output_min")
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
	_, err := fixture.CallFunction(t, "testF32Input_max", float32(math.MaxFloat32))
	if err != nil {
		t.Fatal(err)
	}
}

func TestF32Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testF32Output_max")
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
	_, err := fixture.CallFunction(t, "testF64Input_min", float64(math.SmallestNonzeroFloat64))
	if err != nil {
		t.Fatal(err)
	}
}

func TestF64Output_min(t *testing.T) {
	result, err := fixture.CallFunction(t, "testF64Output_min")
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
	_, err := fixture.CallFunction(t, "testF64Input_max", float64(math.MaxFloat64))
	if err != nil {
		t.Fatal(err)
	}
}

func TestF64Output_max(t *testing.T) {
	result, err := fixture.CallFunction(t, "testF64Output_max")
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
