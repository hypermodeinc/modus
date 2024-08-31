/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"slices"
	"testing"
)

func TestInt8ArrayInput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []int8{0, 1, 2, 3}

	_, err := f.CallFunction("testInt8ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInt8ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt8ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []int16{0, 1, 2, 3}

	_, err := f.CallFunction("testInt16ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInt16ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt16ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []int32{0, 1, 2, 3}

	_, err := f.CallFunction("testInt32ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInt32ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt32ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []int64{0, 1, 2, 3}

	_, err := f.CallFunction("testInt64ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInt64ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testInt64ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []uint8{0, 1, 2, 3}

	_, err := f.CallFunction("testUint8ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUint8ArrayInput_empty(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	_, err := f.CallFunction("testUint8ArrayInput_empty", []uint8{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUint8ArrayInput_null(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	_, err := f.CallFunction("testUint8ArrayInput_null", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUint8ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint8ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint8ArrayOutput_empty")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint8ArrayOutput_null")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestUint16ArrayInput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []uint16{0, 1, 2, 3}

	_, err := f.CallFunction("testUint16ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUint16ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint16ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []uint32{0, 1, 2, 3}

	_, err := f.CallFunction("testUint32ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUint32ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint32ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []uint64{0, 1, 2, 3}

	_, err := f.CallFunction("testUint64ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUint64ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testUint64ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []float32{0, 1, 2, 3}

	_, err := f.CallFunction("testFloat32ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFloat32ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat32ArrayOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	arr := []float64{0, 1, 2, 3}

	_, err := f.CallFunction("testFloat64ArrayInput", arr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFloat64ArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testFloat64ArrayOutput")
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
