/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"testing"
	"time"

	"hypruntime/utils"
)

var testTime, _ = time.Parse(time.RFC3339, "2024-12-31T23:59:59.999999999Z")
var testDuration = time.Duration(5 * time.Second)

func TestTimeInput(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testTimeInput", testTime); err != nil {
		t.Fatal(err)
	}
}

func TestTimePtrInput(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testTimePtrInput", testTime); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testTimePtrInput", &testTime); err != nil {
		t.Fatal(err)
	}
}

func TestTimePtrInput_nil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testTimePtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestTimeOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testTimeOutput")
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
	result, err := fixture.CallFunction(t, "testTimePtrOutput")
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
	result, err := fixture.CallFunction(t, "testTimePtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}

func TestDurationInput(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testDurationInput", testDuration); err != nil {
		t.Fatal(err)
	}
}

func TestDurationPtrInput(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testDurationPtrInput", testDuration); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.CallFunction(t, "testDurationPtrInput", &testDuration); err != nil {
		t.Fatal(err)
	}
}

func TestDurationPtrInput_nil(t *testing.T) {
	if _, err := fixture.CallFunction(t, "testDurationPtrInput_nil", nil); err != nil {
		t.Fatal(err)
	}
}

func TestDurationOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testDurationOutput")
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
	result, err := fixture.CallFunction(t, "testDurationPtrOutput")
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
	result, err := fixture.CallFunction(t, "testDurationPtrOutput_nil")
	if err != nil {
		t.Fatal(err)
	}

	if !utils.HasNil(result) {
		t.Error("expected a nil result")
	}
}
