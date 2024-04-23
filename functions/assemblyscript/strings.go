/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"unicode/utf16"
	"unsafe"

	wasm "github.com/tetratelabs/wazero/api"
)

func ReadString(mem wasm.Memory, offset uint32) (data string, err error) {

	// AssemblyScript managed objects have their classid stored 8 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the class id.
	id, ok := mem.ReadUint32Le(offset - 8)
	if !ok {
		return "", fmt.Errorf("failed to read class id of the WASM object")
	}

	// Make sure the pointer is to a string.
	if id != 2 {
		return "", fmt.Errorf("pointer is not to a string")
	}

	// Read from the buffer and decode it as a string.
	buf, err := readBytes(mem, offset)
	if err != nil {
		return "", err
	}

	return decodeUTF16(buf), nil
}

func WriteString(ctx context.Context, mod wasm.Module, s string) (offset uint32, err error) {
	const classId = 2 // The fixed class id for a string in AssemblyScript.
	bytes := encodeUTF16(s)
	return writeRawBytes(ctx, mod, bytes, classId)
}

func decodeUTF16(bytes []byte) string {

	// Make sure the buffer is valid.
	if len(bytes) == 0 || len(bytes)%2 != 0 {
		return ""
	}

	// Reinterpret []byte as []uint16 to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&bytes[0])
	words := unsafe.Slice((*uint16)(ptr), len(bytes)/2)

	// Decode UTF-16 words to a UTF-8 string.
	str := string(utf16.Decode(words))
	return str
}

func encodeUTF16(str string) []byte {
	if len(str) == 0 {
		return []byte{}
	}
	// Encode the UTF-8 string to UTF-16 words.
	words := utf16.Encode([]rune(str))

	// Reinterpret []uint16 as []byte to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&words[0])
	bytes := unsafe.Slice((*byte)(ptr), len(words)*2)
	return bytes
}
