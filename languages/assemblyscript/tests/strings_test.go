/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"hypruntime/utils"
	"testing"
)

// "Hello World" in Japanese
const testString = "こんにちは、世界"

func TestStringInput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	_, err := f.CallFunction("testStringInput", testString)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStringOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
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
		t.Errorf("expected %q, got %q", testString, r)
	}
}

func TestStringInput_empty(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	_, err := f.CallFunction("testStringInput_empty", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStringOutput_empty(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testStringOutput_empty")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(string); !ok {
		t.Errorf("expected a string, got %T", result)
	} else if r != "" {
		t.Errorf("expected %q, got %q", "", r)
	}
}

func TestNullStringInput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	s := testString

	_, err := f.CallFunction("testNullStringInput", &s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullStringOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testNullStringOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(string); !ok {
		t.Errorf("expected a string, got %T", result)
	} else if r != testString {
		t.Errorf("expected %q, got %q", testString, r)
	}
}

func TestNullStringInput_empty(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	s := ""
	_, err := f.CallFunction("testNullStringInput_empty", &s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullStringOutput_empty(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testNullStringOutput_empty")
	if err != nil {
		t.Fatal(err)
	}

	expected := ""
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(string); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if r != expected {
		t.Errorf("expected %q, got %q", expected, r)
	}
}

func TestNullStringInput_null(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	_, err := f.CallFunction("testNullStringInput_null", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullStringOutput_null(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testNullStringOutput_null")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected no result")
	}
}
