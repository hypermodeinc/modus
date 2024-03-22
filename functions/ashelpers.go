/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf16"
	"unsafe"

	wasm "github.com/tetratelabs/wazero/api"
)

func writeString(ctx context.Context, mod wasm.Module, s string) uint32 {
	buf := encodeUTF16(s)
	return writeObject(ctx, mod, buf, asString)
}

func writeBytes(ctx context.Context, mod wasm.Module, b []byte) uint32 {
	return writeObject(ctx, mod, b, asBytes)
}

func writeObject(ctx context.Context, mod wasm.Module, b []byte, class asClass) uint32 {
	ptr := allocateWasmMemory(ctx, mod, len(b), class)
	mod.Memory().Write(ptr, b)
	return ptr
}

var errPointerIsNotToString = errors.New("pointer is not to a string")

func readString(mem wasm.Memory, offset uint32) (string, error) {

	// AssemblyScript managed objects have their classid stored 8 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the class id.
	id, ok := mem.ReadUint32Le(offset - 8)
	if !ok {
		return "", fmt.Errorf("failed to read class id of the WASM object")
	}

	// Make sure the pointer is to a string.
	if id != uint32(asString) {
		return "", errPointerIsNotToString
	}

	// Read from the buffer and decode it as a string.
	buf, err := readBytes(mem, offset)
	if err != nil {
		return "", err
	}

	return decodeUTF16(buf), nil
}

func readBytes(mem wasm.Memory, offset uint32) ([]byte, error) {

	// The length of AssemblyScript managed objects is stored 4 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the length.
	len, ok := mem.ReadUint32Le(offset - 4)
	if !ok {
		return nil, fmt.Errorf("failed to read buffer length")
	}

	// Handle empty buffers.
	if len == 0 {
		return []byte{}, nil
	}

	// Now read the data into the buffer.
	buf, ok := mem.Read(offset, len)
	if !ok {
		return nil, fmt.Errorf("failed to read buffer data from WASM memory")
	}

	return buf, nil
}

// See https://www.assemblyscript.org/runtime.html#memory-layout
type asClass int64

const (
	asBytes  asClass = 1
	asString asClass = 2
)

func allocateWasmMemory(ctx context.Context, mod wasm.Module, len int, class asClass) uint32 {
	// Allocate a string to hold our buffer within the AssemblyScript module.
	// This uses the `__new` function exported by the AssemblyScript runtime, so it will be garbage collected.
	// See https://www.assemblyscript.org/runtime.html#interface
	newFn := mod.ExportedFunction("__new")
	res, _ := newFn.Call(ctx, uint64(len), uint64(class))
	return uint32(res[0])
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
	// Encode the UTF-8 string to UTF-16 words.
	words := utf16.Encode([]rune(str))

	// Reinterpret []uint16 as []byte to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&words[0])
	bytes := unsafe.Slice((*byte)(ptr), len(words)*2)
	return bytes
}
