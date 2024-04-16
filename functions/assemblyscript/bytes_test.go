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

	buf := []byte{0x01, 0x02, 0x03, 0x04}
	ptr := writeBytes(f.Context, f.Module, buf)
	b, err := readBytes(f.Memory, ptr)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(buf, b) {
		t.Errorf("expected %x, got %x", buf, b)
	}
}
