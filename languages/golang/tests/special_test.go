/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"reflect"
	"testing"
)

func TestMultiOutput(t *testing.T) {
	result, err := fixture.CallFunction(t, "testMultiOutput")
	if err != nil {
		t.Fatal(err)
	}

	expected := []any{123, true, "hello"}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.([]any); !ok {
		t.Errorf("expected []any, got %T", result)
	} else if !reflect.DeepEqual(r, expected) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}
