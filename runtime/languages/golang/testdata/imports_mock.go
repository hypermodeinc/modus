//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

func hostAdd(a, b int) int {
	panic("should not be called")
}

func hostEcho1(message *string) *string {
	panic("should not be called")
}

func hostEcho2(message *string) *string {
	panic("should not be called")
}

func hostEcho3(message *string) *string {
	panic("should not be called")
}

func hostEcho4(message *string) *string {
	panic("should not be called")
}

func hostEncodeStrings1(items *[]string) *string {
	panic("should not be called")
}

func hostEncodeStrings2(items *[]*string) *string {
	panic("should not be called")
}
