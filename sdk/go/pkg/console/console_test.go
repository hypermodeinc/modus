//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package console_test

import (
	"testing"

	"github.com/hypermodeinc/modus/sdk/go/pkg/console"
)

func Test_Log(t *testing.T) {
	msg := "This is a log message."

	console.Log(msg)

	values := console.LogCallStack.Pop()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, but got %d values", len(values))
	}
	if values[0] != "" {
		t.Errorf(`LogCalls[0] = %s; want ""`, values[0])
	}
	if values[1] != msg {
		t.Errorf(`LogCalls[1] = %s; want "%s"`, values[1], msg)
	}
}

func Test_Debug(t *testing.T) {
	msg := "This is a debug message."

	console.Debug(msg)

	values := console.LogCallStack.Pop()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, but got %d values", len(values))
	}
	if values[0] != "debug" {
		t.Errorf(`LogCalls[0] = %s; want "debug"`, values[0])
	}
	if values[1] != msg {
		t.Errorf(`LogCalls[1] = %s; want "%s"`, values[1], msg)
	}
}

func Test_Info(t *testing.T) {
	msg := "This is an info message."

	console.Info(msg)

	values := console.LogCallStack.Pop()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, but got %d values", len(values))
	}
	if values[0] != "info" {
		t.Errorf(`LogCalls[0] = %s; want "info"`, values[0])
	}
	if values[1] != msg {
		t.Errorf(`LogCalls[1] = %s; want "%s"`, values[1], msg)
	}
}

func Test_Warn(t *testing.T) {
	msg := "This is a warning message."

	console.Warn(msg)

	values := console.LogCallStack.Pop()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, but got %d values", len(values))
	}
	if values[0] != "warning" {
		t.Errorf(`LogCalls[0] = %s; want "warning"`, values[0])
	}
	if values[1] != msg {
		t.Errorf(`LogCalls[1] = %s; want "%s"`, values[1], msg)
	}
}

func Test_Error(t *testing.T) {
	msg := "This is an error message."

	console.Error(msg)

	values := console.LogCallStack.Pop()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, but got %d values", len(values))
	}
	if values[0] != "error" {
		t.Errorf(`LogCalls[0] = %s; want "error"`, values[0])
	}
	if values[1] != msg {
		t.Errorf(`LogCalls[1] = %s; want "%s"`, values[1], msg)
	}
}
func Test_Assert(t *testing.T) {
	// Test case 1: condition is true
	sizeBefore := console.LogCallStack.Size()
	console.Assert(true, "Condition is true")
	sizeAfter := console.LogCallStack.Size()
	if sizeBefore != sizeAfter {
		values := console.LogCallStack.Pop()
		t.Errorf("Expected no values logged, but got %v", values)
	}

	// Test case 2: condition is false
	console.Assert(false, "Condition is false")
	values := console.LogCallStack.Pop()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, but got %d values", len(values))
	}
	if values[0] != "error" {
		t.Errorf(`LogCalls[0] = %s; want "error"`, values[0])
	}
	if values[1] != "Assertion failed: Condition is false" {
		t.Errorf(`LogCalls[1] = %s; want "Assertion failed: Condition is false"`, values[1])
	}
}
