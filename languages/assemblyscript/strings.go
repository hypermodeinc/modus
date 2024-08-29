/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"

	"hmruntime/utils"
)

func (wa *wasmAdapter) readString(offset uint32) (data string, err error) {
	if offset == 0 {
		return "", nil
	}

	// AssemblyScript managed objects have their classid stored 8 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the class id.
	id, ok := wa.mod.Memory().ReadUint32Le(offset - 8)
	if !ok {
		return "", fmt.Errorf("failed to read class id of the WASM object")
	}

	// Make sure the pointer is to a string.
	if id != 2 {
		return "", fmt.Errorf("pointer is not to a string")
	}

	// Read from the buffer and decode it as a string.
	buf, err := wa.readBytes(offset)
	if err != nil {
		return "", err
	}

	return utils.DecodeUTF16(buf), nil
}

func (wa *wasmAdapter) writeString(ctx context.Context, s string) (offset uint32, err error) {
	const classId = 2 // The fixed class id for a string in AssemblyScript.
	bytes := utils.EncodeUTF16(s)
	return wa.writeRawBytes(ctx, bytes, classId)
}
