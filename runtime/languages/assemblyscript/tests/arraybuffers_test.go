/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"bytes"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

func TestArrayBufferInput(t *testing.T) {
	arr := []byte{1, 2, 3, 4}
	fnName := "testArrayBufferInput"

	if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
	if arr, ok := utils.ConvertToSliceOf[any](arr); !ok {
		t.Error("failed conversion to interface slice")
	} else if _, err := fixture.CallFunction(t, fnName, arr); err != nil {
		t.Error(err)
	}
}

func TestArrayBufferOutput(t *testing.T) {
	fnName := "testArrayBufferOutput"
	result, err := fixture.CallFunction(t, fnName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{1, 2, 3, 4}
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]byte); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if !bytes.Equal(expected, r) {
		t.Errorf("expected %x, got %x", expected, r)
	}
}
