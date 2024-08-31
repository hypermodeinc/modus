/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"testing"

	"hypruntime/utils"
)

// "Hello World" in Japanese
const testString = "こんにちは、世界"

func TestStringInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testStringInput", testString); err != nil {
		t.Fatal(err)
	}
}

func TestStringPtrInput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	s := testString
	if _, err := f.CallFunction("testStringPtrInput", s); err != nil {
		t.Fatal(err)
	}
	if _, err := f.CallFunction("testStringPtrInput", &s); err != nil {
		t.Fatal(err)
	}
}

func TestStringPtrInput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	if _, err := f.CallFunction("testStringPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestStringOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testStringOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(string); !ok {
		t.Errorf("expected a string, got %T", result)
	} else if r != testString {
		t.Errorf("expected %s, got %s", testString, r)
	}
}

func TestStringPtrOutput(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testStringPtrOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*string); !ok {
		t.Errorf("expected a *string, got %T", result)
	} else if *r != testString {
		t.Errorf("expected %s, got %s", testString, *r)
	}
}

func TestStringPtrOutput_nil(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testStringPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
