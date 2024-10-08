/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript_test

import (
	"testing"
	"time"
)

func Test_DateNow(t *testing.T) {
	fnName := "now"
	startTs := time.Now().UnixMilli()
	time.Sleep(1 * time.Millisecond)

	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	resultTs := result.(int64)
	if resultTs <= startTs {
		t.Errorf("expected timestamp greater than %d, got %d", startTs, resultTs)
	}
}

func Test_PerformanceNow(t *testing.T) {
	fnName := "spin"
	const expectedDuration = 100 // ms

	result, err := fixture.CallFunction(t, fnName, expectedDuration)
	if err != nil {
		t.Fatal(err)
	}

	actualDuration := result.(int64)
	if actualDuration < expectedDuration {
		t.Errorf("expected duration of at least than %d, got %d", expectedDuration, actualDuration)
	}
}
