/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import "time"

var testTime, _ = time.Parse(time.RFC3339, "2024-12-31T23:59:59.999999999Z")
var testDuration = time.Duration(5 * time.Second)

func TestTimeInput(t time.Time) {
	assertEqual(testTime, t)
}

func TestTimePtrInput(t *time.Time) {
	assertEqual(testTime, *t)
}

func TestTimePtrInput_nil(t *time.Time) {
	assertNil(t)
}

func TestTimeOutput() time.Time {
	return testTime
}

func TestTimePtrOutput() *time.Time {
	return &testTime
}

func TestTimePtrOutput_nil() *time.Time {
	return nil
}

func TestDurationInput(d time.Duration) {
	assertEqual(testDuration, d)
}

func TestDurationPtrInput(d *time.Duration) {
	assertEqual(testDuration, *d)
}

func TestDurationPtrInput_nil(d *time.Duration) {
	assertNil(d)
}

func TestDurationOutput() time.Duration {
	return testDuration
}

func TestDurationPtrOutput() *time.Duration {
	return &testDuration
}

func TestDurationPtrOutput_nil() *time.Duration {
	return nil
}
