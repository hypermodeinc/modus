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
	"bytes"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

// "Hello World" in Japanese
const testString = "こんにちは、世界"

var testUTF16 = []byte{
	0x53, 0x30, 0x93, 0x30,
	0x6b, 0x30, 0x61, 0x30,
	0x6f, 0x30, 0x01, 0x30,
	0x16, 0x4e, 0x4c, 0x75,
}

func Test_EncodeUTF16(t *testing.T) {

	arr := utils.EncodeUTF16(testString)

	if !bytes.Equal(arr, testUTF16) {
		t.Errorf("expected %x, got %x", testUTF16, arr)
	}
}

func Test_DecodeUTF16(t *testing.T) {

	str := utils.DecodeUTF16(testUTF16)

	if str != testString {
		t.Errorf("expected %s, got %s", testString, str)
	}
}
