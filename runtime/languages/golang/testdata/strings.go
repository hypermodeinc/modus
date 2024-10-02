/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

// "Hello World" in Japanese
const testString = "こんにちは、世界"

func TestStringInput(s string) {
	assertEqual(testString, s)
}

func TestStringPtrInput(s *string) {
	assertEqual(testString, *s)
}

func TestStringPtrInput_nil(s *string) {
	assertNil(s)
}

func TestStringOutput() string {
	return testString
}

func TestStringPtrOutput() *string {
	s := testString
	return &s
}

func TestStringPtrOutput_nil() *string {
	return nil
}
