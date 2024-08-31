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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	_, err := f.CallFunction("testDateInput", testTime)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDateOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testDateOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	dt := testTime

	_, err := f.CallFunction("testNullDateInput", &dt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullDateOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testNullDateOutput")
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
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	_, err := f.CallFunction("testNullDateInput_null", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNullDateOutput_null(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.CallFunction("testNullDateOutput_null")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected no result")
	}
}
