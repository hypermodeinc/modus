/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils_test

import (
	"bytes"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/stretchr/testify/assert"
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

type mockOutputBuffers struct {
	stdOut *bytes.Buffer
	stdErr *bytes.Buffer
}

func (m mockOutputBuffers) StdOut() *bytes.Buffer {
	return m.stdOut
}

func (m mockOutputBuffers) StdErr() *bytes.Buffer {
	return m.stdErr
}

func TestTransformConsoleOutput(t *testing.T) {
	tests := []struct {
		stdOut   string
		stdErr   string
		expected []utils.LogMessage
	}{
		{
			stdOut:   "",
			stdErr:   "",
			expected: []utils.LogMessage{},
		},
		{
			stdOut: "Info: This is an info message\n",
			stdErr: "Error: This is an error message\n",
			expected: []utils.LogMessage{
				{Level: "info", Message: "This is an info message"},
				{Level: "error", Message: "This is an error message"},
			},
		},
		{
			stdOut: "Debug: This is a debug message\nWarning: This is a warning message\n",
			stdErr: "panic: This is a fatal message\n",
			expected: []utils.LogMessage{
				{Level: "debug", Message: "This is a debug message"},
				{Level: "warning", Message: "This is a warning message"},
				{Level: "fatal", Message: "This is a fatal message"},
			},
		},
	}

	for _, test := range tests {
		buffers := mockOutputBuffers{
			stdOut: bytes.NewBufferString(test.stdOut),
			stdErr: bytes.NewBufferString(test.stdErr),
		}
		result := utils.TransformConsoleOutput(buffers)
		assert.Equal(t, test.expected, result)
	}
}
