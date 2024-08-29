/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"testing"
	"time"
)

func Test_DateNow(t *testing.T) {
	t.Parallel()

	startTs := time.Now().UnixMilli()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	fn := f.Module.ExportedFunction("now")
	res, err := fn.Call(f.Context)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 {
		t.Errorf("expected 1 result, got %d", len(res))
	}

	resultTs := int64(res[0])

	if resultTs <= startTs {
		t.Errorf("expected timestamp greater than %d, got %d", startTs, resultTs)
	}
}

func Test_PerformanceNow(t *testing.T) {
	t.Parallel()

	const expectedDuration = 100 // ms

	f := NewASWasmTestFixture(t)
	defer f.Close()

	fn := f.Module.ExportedFunction("spin")
	res, err := fn.Call(f.Context, expectedDuration)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 {
		t.Errorf("expected 1 result, got %d", len(res))
	}

	actualDuration := int64(res[0])

	if actualDuration < expectedDuration {
		t.Errorf("expected duration of at least than %d, got %d", expectedDuration, actualDuration)
	}
}
