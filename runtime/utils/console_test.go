/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils_test

import (
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

func TestIsError(t *testing.T) {
	tests := []struct {
		logMessage utils.LogMessage
		expected   bool
	}{
		{utils.LogMessage{Level: "debug", Message: "This is a debug message"}, false},
		{utils.LogMessage{Level: "info", Message: "This is an info message"}, false},
		{utils.LogMessage{Level: "warning", Message: "This is a warning message"}, false},
		{utils.LogMessage{Level: "error", Message: "This is an error message"}, true},
		{utils.LogMessage{Level: "fatal", Message: "This is a fatal message"}, true},
	}

	for _, test := range tests {
		result := test.logMessage.IsError()
		if result != test.expected {
			t.Errorf("For log message '%v', expected IsError to be '%v', but got '%v'", test.logMessage, test.expected, result)
		}
	}
}

func TestSplitConsoleOutputLine(t *testing.T) {
	tests := []struct {
		input           string
		expectedLevel   string
		expectedMessage string
	}{
		{"Debug: This is a debug message", "debug", "This is a debug message"},
		{"Info: This is an info message", "info", "This is an info message"},
		{"Warning: This is a warning message", "warning", "This is a warning message"},
		{"Error: This is an error message", "error", "This is an error message"},
		{"abort: This is a fatal message", "fatal", "This is a fatal message"},
		{"panic: This is another fatal message", "fatal", "This is another fatal message"},
		{"This is a message without level", "", "This is a message without level"},
	}

	for _, test := range tests {
		level, message := utils.SplitConsoleOutputLine(test.input)
		if level != test.expectedLevel || message != test.expectedMessage {
			t.Errorf("For input '%s', expected level '%s' and message '%s', but got level '%s' and message '%s'", test.input, test.expectedLevel, test.expectedMessage, level, message)
		}
	}
}
