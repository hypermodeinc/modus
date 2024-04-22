/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"errors"
	"fmt"

	wasm "github.com/tetratelabs/wazero/api"
)

func readBytes(mem wasm.Memory, offset uint32) ([]byte, error) {

	// The length of AssemblyScript managed objects is stored 4 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the length.
	len, ok := mem.ReadUint32Le(offset - 4)
	if !ok {
		return nil, errors.New("failed to read buffer length")
	}

	// Handle empty buffers.
	if len == 0 {
		return []byte{}, nil
	}

	// Now read the data into the buffer.
	buf, ok := mem.Read(offset, len)
	if !ok {
		return nil, errors.New("failed to read buffer data from WASM memory")
	}

	return buf, nil
}

func writeBytes(ctx context.Context, mod wasm.Module, bytes []byte) (uint32, error) {
	const classId = 1 // The fixed class id for an ArrayBuffer in AssemblyScript.
	return writeRawBytes(ctx, mod, bytes, classId)
}

func writeRawBytes(ctx context.Context, mod wasm.Module, bytes []byte, classId uint32) (uint32, error) {
	size := len(bytes)
	offset, err := allocateWasmMemory(ctx, mod, size, classId)
	if err != nil {
		return 0, err
	}

	ok := mod.Memory().Write(offset, bytes)
	if !ok {
		return 0, fmt.Errorf("failed to write to WASM memory (%d bytes, class id %d)", size, classId)
	}

	return offset, nil
}
