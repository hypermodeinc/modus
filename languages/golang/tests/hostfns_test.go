/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"testing"

	"hypruntime/wasmhost"
)

func getTestHostFunctionRegistrations() []func(*wasmhost.WasmHost) error {
	return []func(*wasmhost.WasmHost) error{
		func(host *wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "add", hostAdd)
		},
		func(host *wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "echo", hostEcho)
		},
	}
}

func hostAdd(a, b int) int {
	return a + b
}

func hostEcho(s *string) string {
	return "echo: " + *s
}

func TestHostFn_add(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("add", 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int); !ok {
		t.Errorf("expected an int, got %T", result)
	} else if r != 3 {
		t.Errorf("expected %d, got %d", 3, r)
	}
}

func TestHostFn_echo(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("echo", "hello")
	if err != nil {
		t.Fatal(err)
	}

	expected := "echo: hello"
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(string); !ok {
		t.Errorf("expected a string, got %T", result)
	} else if r != expected {
		t.Errorf("expected %s, got %s", expected, r)
	}
}
