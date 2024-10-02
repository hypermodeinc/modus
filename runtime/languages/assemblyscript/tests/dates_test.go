/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"hypruntime/utils"
	"testing"
	"time"
)

var testTime, _ = time.Parse(time.RFC3339, "2024-12-31T23:59:59.999Z")

func TestDateInput(t *testing.T) {
	fnName := "testDateInput"
	if _, err := fixture.CallFunction(t, fnName, testTime); err != nil {
		t.Error(err)
	}
}

func TestDateOutput(t *testing.T) {
	fnName := "testDateOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(time.Time); !ok {
		t.Errorf("expected %T, got %T", testTime, result)
	} else if r != testTime {
		t.Errorf("expected %q, got %q", testTime, r)
	}
}

func TestNullDateInput(t *testing.T) {
	fnName := "testNullDateInput"
	dt := testTime

	if _, err := fixture.CallFunction(t, fnName, &dt); err != nil {
		t.Error(err)
	}
}

func TestNullDateOutput(t *testing.T) {
	fnName := "testNullDateOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(time.Time); !ok {
		t.Errorf("expected a %T, got %T", testTime, result)
	} else if r != testTime {
		t.Errorf("expected %q, got %q", testTime, r)
	}
}

func TestNullDateInput_null(t *testing.T) {
	fnName := "testNullDateInput_null"
	if _, err := fixture.CallFunction(t, fnName, nil); err != nil {
		t.Error(err)
	}
}

func TestNullDateOutput_null(t *testing.T) {
	fnName := "testNullDateOutput_null"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected no result")
	}
}
