/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"bytes"
	"testing"

	"hmruntime/testutils"
)

func Test_ReadWriteBuffer(t *testing.T) {
	f := testutils.NewWasmTestFixture()
	defer f.Close()

	wa := _wasmAdapter

	buf := []byte{0x01, 0x02, 0x03, 0x04}
	ptr, err := wa.writeBytes(f.Context, f.Module, buf)
	if err != nil {
		t.Error(err)
	}

	b, err := wa.readBytes(f.Memory, ptr)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(buf, b) {
		t.Errorf("expected %x, got %x", buf, b)
	}
}
