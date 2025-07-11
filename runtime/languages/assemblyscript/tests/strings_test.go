/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript_test

import (
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

// "Hello World" in Japanese
const testString = "こんにちは、世界"

func TestStringInput(t *testing.T) {
	fnName := "testStringInput"
	if _, err := fixture.CallFunction(t, fnName, testString); err != nil {
		t.Error(err)
	}
}

func TestStringOutput(t *testing.T) {
	fnName := "testStringOutput"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testStringInput_empty"
	if _, err := fixture.CallFunction(t, fnName, ""); err != nil {
		t.Error(err)
	}
}

func TestStringOutput_empty(t *testing.T) {
	fnName := "testStringOutput_empty"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testNullStringInput"
	s := testString

	if _, err := fixture.CallFunction(t, fnName, &s); err != nil {
		t.Error(err)
	}
}

func TestNullStringOutput(t *testing.T) {
	fnName := "testNullStringOutput"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testNullStringInput_empty"
	s := ""

	if _, err := fixture.CallFunction(t, fnName, &s); err != nil {
		t.Error(err)
	}
}

func TestNullStringOutput_empty(t *testing.T) {
	fnName := "testNullStringOutput_empty"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testNullStringInput_null"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestNullStringOutput_null(t *testing.T) {
	fnName := "testNullStringOutput_null"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected no result")
	}
}
