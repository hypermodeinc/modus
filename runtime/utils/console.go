/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"bytes"
	"strings"
)

type LogMessage struct {
	Level   string `json:"level,omitempty"`
	Message string `json:"message"`
}

func (l LogMessage) IsError() bool {
	return l.Level == "error" || l.Level == "fatal"
}

func TransformConsoleOutput(buffers OutputBuffers) []LogMessage {
	return append(transformConsoleOutputLines(buffers.StdOut()), transformConsoleOutputLines(buffers.StdErr())...)
}

func transformConsoleOutputLines(buf *bytes.Buffer) []LogMessage {
	lines := strings.Split(buf.String(), "\n")
	messages := make([]LogMessage, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			level, message := SplitConsoleOutputLine(line)
			messages = append(messages, LogMessage{level, message})
		}
	}
	return messages
}

func SplitConsoleOutputLine(line string) (level string, message string) {
	a := strings.SplitAfterN(line, ": ", 2)
	if len(a) == 2 {
		switch a[0] {
		case "Debug: ":
			return "debug", a[1]
		case "Info: ":
			return "info", a[1]
		case "Warning: ":
			return "warning", a[1]
		case "Error: ":
			return "error", a[1]
		case "abort: ", "panic: ":
			return "fatal", a[1]
		}
	}
	return "", line
}
