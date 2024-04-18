/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"

	wasm "github.com/tetratelabs/wazero/api"
)

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

func writeBytes(ctx context.Context, mod wasm.Module, bytes []byte) uint32 {
	return writeObject(ctx, mod, bytes, 1)
}
