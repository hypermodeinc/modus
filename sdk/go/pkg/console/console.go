/*
 * Copyright 2024 Hypermode, Inc.
 */

package console

import "fmt"

func Assert(condition bool, message string) {
	if !condition {
		Error("Assertion failed: " + message)
	}
}

func Log(message string) {
	log("", message)
}

func Logf(format string, args ...any) {
	Log(fmt.Sprintf(format, args...))
}

func Debug(message string) {
	log("debug", message)
}

func Debugf(format string, args ...any) {
	Debug(fmt.Sprintf(format, args...))
}

func Info(message string) {
	log("info", message)
}

func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

func Warn(message string) {
	log("warning", message)
}

func Warnf(format string, args ...any) {
	Warn(fmt.Sprintf(format, args...))
}

func Error(message string) {
	log("error", message)
}

func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}
