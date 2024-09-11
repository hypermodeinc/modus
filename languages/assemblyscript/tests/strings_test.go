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
	_, err := fixture.CallFunction(t, "testStringInput", testString)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStringOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStringOutput")
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
	_, err := fixture.CallFunction(t, "testStringInput_empty", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStringOutput_empty(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStringOutput_empty")
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
	s := testString

	_, err := fixture.CallFunction(t, "testNullStringInput", &s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullStringOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testNullStringOutput")
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
	s := ""
	_, err := fixture.CallFunction(t, "testNullStringInput_empty", &s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullStringOutput_empty(t *testing.T) {
	result, err := fixture.CallFunction(t, "testNullStringOutput_empty")
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
	_, err := fixture.CallFunction(t, "testNullStringInput_null", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullStringOutput_null(t *testing.T) {
	result, err := fixture.CallFunction(t, "testNullStringOutput_null")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected no result")
	}
}
