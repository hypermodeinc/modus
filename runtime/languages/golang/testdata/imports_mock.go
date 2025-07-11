//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

func hostAdd(a, b int32) int32 {
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
