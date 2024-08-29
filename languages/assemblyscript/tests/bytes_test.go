/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"bytes"
	"testing"
)

var testBytes = []byte{0x01, 0x02, 0x03, 0x04}

func TestArrayBufferInput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	_, err := f.InvokeFunction("testArrayBufferInput", testBytes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayBufferOutput(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("testArrayBufferOutput")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]byte); !ok {
		t.Errorf("expected a []byte, got %T", result)
	} else if !bytes.Equal(testBytes, r) {
		t.Errorf("expected %x, got %x", testBytes, r)
	}
}
