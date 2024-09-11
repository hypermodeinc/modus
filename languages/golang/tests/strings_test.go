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
	if _, err := fixture.CallFunction(t, "testStringInput", testString); err != nil {
		t.Fatal(err)
	}
}

func TestStringPtrInput(t *testing.T) {
	s := testString
	if _, err := fixture.CallFunction(t, "testStringPtrInput", s); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testStringPtrInput", &s); err != nil {
		t.Fatal(err)
	}
}

func TestStringPtrInput_nil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testStringPtrInput_nil", nil); err != nil {
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
		t.Errorf("expected %s, got %s", testString, r)
	}
}

func TestStringPtrOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testStringPtrOutput")
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
	result, err := fixture.CallFunction(t, "testStringPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
