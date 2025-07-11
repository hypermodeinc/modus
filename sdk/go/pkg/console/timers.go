/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package console

import (
	"strconv"
	"time"
)

var timers = make(map[string]time.Time)

func Time(label ...string) {
	now := time.Now()

	l, ok := getLabel(label)
	if !ok {
		Error("Invalid number of arguments passed to console.Time function.")
		return
	}

	if _, ok := timers[l]; ok {
		Warn("Timer label '" + l + "' already exists.")
		return
	}
	timers[l] = now
}

func TimeLog(label ...string) {
	now := time.Now()

	l, ok := getLabel(label)
	if !ok {
		Error("Invalid number of arguments passed to console.TimeLog function.")
		return
	}

	t, ok := timers[l]
	if !ok {
		Warn("Timer label '" + l + "' does not exist.")
		return
	}
	ms := now.Sub(t).Milliseconds()
	Info(l + ": " + strconv.FormatInt(ms, 10) + "ms")
}

func TimeEnd(label ...string) {
	now := time.Now()

	l, ok := getLabel(label)
	if !ok {
		Error("Invalid number of arguments passed to console.TimeEnd function.")
		return
	}

	t, ok := timers[l]
	if !ok {
		Warn("Timer label '" + l + "' does not exist.")
		return
	}
	ms := now.Sub(t).Milliseconds()
	Info(l + ": " + strconv.FormatInt(ms, 10) + "ms")
	delete(timers, l)
}

func getLabel(labels []string) (string, bool) {
	switch len(labels) {
	case 0:
		return "default", true
	case 1:
		return labels[0], true
	default:
		return "", false
	}
}
