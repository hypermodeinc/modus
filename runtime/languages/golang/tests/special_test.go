/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang_test

import (
	"reflect"
	"testing"
)

func TestMultiOutput(t *testing.T) {
	fnName := "testMultiOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []any{123, true, "hello"}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]any); !ok {
		t.Errorf("expected []any, got %T", result)
	} else if !reflect.DeepEqual(r, expected) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}
