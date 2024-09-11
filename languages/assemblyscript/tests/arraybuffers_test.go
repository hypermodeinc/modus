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
	_, err := fixture.CallFunction(t, "testArrayBufferInput", testBytes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayBufferOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testArrayBufferOutput")
	if err != nil {
		t.Fatal(err)
	}

	expected := testBytes
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]byte); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !bytes.Equal(expected, r) {
		t.Errorf("expected %x, got %x", expected, r)
	}
}
