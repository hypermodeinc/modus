/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"bytes"
	"testing"

	"hmruntime/testutils"
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

	arr := encodeUTF16(testString)

	if !bytes.Equal(arr, testUTF16) {
		t.Errorf("expected %x, got %x", testUTF16, arr)
	}
}

func Test_DecodeUTF16(t *testing.T) {

	str := decodeUTF16(testUTF16)

	if str != testString {
		t.Errorf("expected %s, got %s", testString, str)
	}
}

func Test_ReadWriteString(t *testing.T) {
	f := testutils.NewWasmTestFixture()
	defer f.Close()

	ptr, err := WriteString(f.Context, f.Module, testString)
	if err != nil {
		t.Error(err)
	}

	str, err := ReadString(f.Memory, ptr)
	if err != nil {
		t.Error(err)
	}

	if str != testString {
		t.Errorf("expected %s, got %s", testString, str)
	}
}
