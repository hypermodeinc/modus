/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package console

import "fmt"

func Assert(condition bool, message string) {
	if !condition {
		Error("Assertion failed: " + message)
	}
}

func Log(message string) {
	hostLogMessage("", message)
}

func Logf(format string, args ...any) {
	Log(fmt.Sprintf(format, args...))
}

func Debug(message string) {
	hostLogMessage("debug", message)
}

func Debugf(format string, args ...any) {
	Debug(fmt.Sprintf(format, args...))
}

func Info(message string) {
	hostLogMessage("info", message)
}

func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

func Warn(message string) {
	hostLogMessage("warning", message)
}

func Warnf(format string, args ...any) {
	Warn(fmt.Sprintf(format, args...))
}

func Error(message string) {
	hostLogMessage("error", message)
}

func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}
