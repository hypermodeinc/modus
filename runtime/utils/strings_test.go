/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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

func Test_SanitizeUTF8(t *testing.T) {
	// Test with a valid UTF-8 string
	validUTF8 := []byte("Hello, 世界")
	sanitized := utils.SanitizeUTF8(validUTF8)
	if !bytes.Equal(sanitized, validUTF8) {
		t.Errorf("expected %s, got %s", validUTF8, sanitized)
	}

	// Test with an invalid UTF-8 sequence
	invalidUTF8 := []byte{0xff, 0xfe, 0xfd}
	sanitized = utils.SanitizeUTF8(invalidUTF8)
	if len(sanitized) != 0 {
		t.Errorf("expected empty slice for invalid UTF-8, got %s", sanitized)
	}

	// Test with a mix of valid and invalid UTF-8
	mixedUTF8 := []byte("Hello\xffWorld")
	sanitized = utils.SanitizeUTF8(mixedUTF8)
	expected := []byte("HelloWorld")
	if !bytes.Equal(sanitized, expected) {
		t.Errorf("expected %s, got %s", expected, sanitized)
	}

	// Test with some null bytes
	nullBytes := []byte("Hello\x00World")
	sanitized = utils.SanitizeUTF8(nullBytes)
	expected = []byte("HelloWorld")
	if !bytes.Equal(sanitized, expected) {
		t.Errorf("expected %s, got %s", expected, sanitized)
	}
}
