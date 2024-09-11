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
	_, err := fixture.CallFunction(t, "testDateInput", testTime)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDateOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testDateOutput")
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
	dt := testTime

	_, err := fixture.CallFunction(t, "testNullDateInput", &dt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullDateOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testNullDateOutput")
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
	_, err := fixture.CallFunction(t, "testNullDateInput_null", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullDateOutput_null(t *testing.T) {
	result, err := fixture.CallFunction(t, "testNullDateOutput_null")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected no result")
	}
}
