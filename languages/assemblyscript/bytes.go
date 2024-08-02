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

func (wa *wasmAdapter) readBytes(mem wasm.Memory, offset uint32) (data []byte, err error) {

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

func (wa *wasmAdapter) writeBytes(ctx context.Context, mod wasm.Module, bytes []byte) (offset uint32, err error) {
	const classId = 1 // The fixed class id for an ArrayBuffer in AssemblyScript.
	return wa.writeRawBytes(ctx, mod, bytes, classId)
}

func (wa *wasmAdapter) writeRawBytes(ctx context.Context, mod wasm.Module, bytes []byte, classId uint32) (offset uint32, err error) {
	size := uint32(len(bytes))
	offset, err = wa.allocateWasmMemory(ctx, mod, size, classId)
	if err != nil {
		return 0, err
	}

	ok := mod.Memory().Write(offset, bytes)
	if !ok {
		return 0, fmt.Errorf("failed to write to WASM memory (%d bytes, class id %d)", size, classId)
	}

	return offset, nil
}
