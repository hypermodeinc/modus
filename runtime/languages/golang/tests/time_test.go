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
	"time"

	"github.com/hypermodeinc/modus/runtime/utils"
)

var testTime, _ = time.Parse(time.RFC3339, "2024-12-31T23:59:59.999999999Z")
var testDuration = time.Duration(5 * time.Second)

func TestTimeInput(t *testing.T) {
	fnName := "testTimeInput"
	if _, err := fixture.CallFunction(t, fnName, testTime); err != nil {
		t.Error(err)
	}
}

func TestTimePtrInput(t *testing.T) {
	fnName := "testTimePtrInput"
	if _, err := fixture.CallFunction(t, fnName, testTime); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &testTime); err != nil {
		t.Error(err)
	}
}

func TestTimePtrInput_nil(t *testing.T) {
	fnName := "testTimePtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestTimeOutput(t *testing.T) {
	fnName := "testTimeOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(time.Time); !ok {
		t.Errorf("expected a time.Time, got %T", result)
	} else if r != testTime {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestTimePtrOutput(t *testing.T) {
	fnName := "testTimePtrOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*time.Time); !ok {
		t.Errorf("expected a *time.Time, got %T", result)
	} else if *r != testTime {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestTimePtrOutput_nil(t *testing.T) {
	fnName := "testTimePtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestDurationInput(t *testing.T) {
	fnName := "testDurationInput"
	if _, err := fixture.CallFunction(t, fnName, testDuration); err != nil {
		t.Error(err)
	}
}

func TestDurationPtrInput(t *testing.T) {
	fnName := "testDurationPtrInput"
	if _, err := fixture.CallFunction(t, fnName, testDuration); err != nil {
		t.Error(err)
	}
	if _, err := fixture.CallFunction(t, fnName, &testDuration); err != nil {
		t.Error(err)
	}
}

func TestDurationPtrInput_nil(t *testing.T) {
	fnName := "testDurationPtrInput_nil"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestDurationOutput(t *testing.T) {
	fnName := "testDurationOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(time.Duration); !ok {
		t.Errorf("expected a time.Duration, got %T", result)
	} else if r != testDuration {
		t.Errorf("expected %v, got %v", true, r)
	}
}

func TestDurationPtrOutput(t *testing.T) {
	fnName := "testDurationPtrOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(*time.Duration); !ok {
		t.Errorf("expected a *time.Duration, got %T", result)
	} else if *r != testDuration {
		t.Errorf("expected %v, got %v", true, *r)
	}
}

func TestDurationPtrOutput_nil(t *testing.T) {
	fnName := "testDurationPtrOutput_nil"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
