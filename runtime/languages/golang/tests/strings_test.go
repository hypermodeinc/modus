/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang_test

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

func TestStringPtrInput(t *testing.T) {
	fnName := "testStringPtrInput"
	s := testString

	if _, err := fixture.CallFunction(t, fnName, s); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &s); err != nil {
		t.Error(err)
	}
}

func TestStringPtrInput_nil(t *testing.T) {
	fnName := "testStringPtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
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
		t.Errorf("expected %s, got %s", testString, r)
	}
}

func TestStringPtrOutput(t *testing.T) {
	fnName := "testStringPtrOutput"
	result, err := fixture.CallFunction(t, fnName)
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
	fnName := "testStringPtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
