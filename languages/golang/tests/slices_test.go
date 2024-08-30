/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"bytes"
	"slices"
	"testing"
)

func TestSliceInput_byte(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = []byte{0x01, 0x02, 0x03, 0x04}

	if _, err := f.InvokeFunction("testSliceInput_byte", val); err != nil {
		t.Fatal(err)
	}
}

func TestSliceInput_intPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getIntPtrSlice()

	if _, err := f.InvokeFunction("testSliceInput_intPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestSliceInput_string(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = []string{"abc", "def", "ghi"}

	if _, err := f.InvokeFunction("testSliceInput_string", val); err != nil {
		t.Fatal(err)
	}
}

func TestSliceInput_stringPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	var val = getStringPtrSlice()

	if _, err := f.InvokeFunction("testSliceInput_stringPtr", val); err != nil {
		t.Fatal(err)
	}
}

func TestSliceOutput_byte(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testSliceOutput_byte")
	if err != nil {
		t.Fatal(err)
	}

	var expected = []byte{0x01, 0x02, 0x03, 0x04}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]byte); !ok {
		t.Errorf("expected a []byte, got %T", result)
	} else if !bytes.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestSliceOutput_intPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testSliceOutput_intPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getIntPtrSlice()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]*int); !ok {
		t.Errorf("expected a []*int, got %T", result)
	} else if !slices.EqualFunc(expected, r, func(a, b *int) bool { return *a == *b }) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestSliceOutput_string(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testSliceOutput_string")
	if err != nil {
		t.Fatal(err)
	}

	var expected = []string{"abc", "def", "ghi"}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]string); !ok {
		t.Errorf("expected a []string, got %T", result)
	} else if !slices.Equal(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func TestSliceOutput_stringPtr(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testSliceOutput_stringPtr")
	if err != nil {
		t.Fatal(err)
	}

	var expected = getStringPtrSlice()
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]*string); !ok {
		t.Errorf("expected a []*string, got %T", result)
	} else if !slices.EqualFunc(expected, r, func(a, b *string) bool { return *a == *b }) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}

func getIntPtrSlice() []*int {
	a := 11
	b := 22
	c := 33
	return []*int{&a, &b, &c}
}

func getStringPtrSlice() []*string {
	a := "abc"
	b := "def"
	c := "ghi"
	return []*string{&a, &b, &c}
}
