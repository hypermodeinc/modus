/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"context"
	"testing"

	"hypruntime/hostfunctions"
	"hypruntime/testutils"
	"hypruntime/utils"
	"hypruntime/wasmhost"
)

func getTestHostFunctionRegistrations() []func(wasmhost.WasmHost) error {
	return []func(wasmhost.WasmHost) error{
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("hypermode", "log", hostLog)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "add", hostAdd)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "echo1", hostEcho1)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "echo2", hostEcho2)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "echo3", hostEcho3)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "echo4", hostEcho4)
		},
	}
}

func hostLog(ctx context.Context, level, message string) {
	if utils.HypermodeDebugEnabled() {
		hostfunctions.LogFunctionMessage(ctx, level, message)
	}
	t := testutils.GetTestT(ctx)
	t.Logf("[%s] %s", level, message)
}

func hostAdd(a, b int) int {
	return a + b
}

func hostEcho1(s string) string {
	return "echo: " + s
}

func hostEcho2(s *string) string {
	return "echo: " + *s
}

func hostEcho3(s string) *string {
	result := "echo: " + s
	return &result
}

func hostEcho4(s *string) *string {
	result := "echo: " + *s
	return &result
}

func TestHostFn_add(t *testing.T) {
	result, err := fixture.CallFunction(t, "add", 1, 2)
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

func TestHostFn_echo1_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "echo1", "hello")
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

func TestHostFn_echo1_stringPtr(t *testing.T) {
	s := "hello"
	result, err := fixture.CallFunction(t, "echo1", &s)
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

func TestHostFn_echo2_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "echo2", "hello")
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

func TestHostFn_echo2_stringPtr(t *testing.T) {
	s := "hello"
	result, err := fixture.CallFunction(t, "echo2", &s)
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

func TestHostFn_echo3_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "echo3", "hello")
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

func TestHostFn_echo3_stringPtr(t *testing.T) {
	s := "hello"
	result, err := fixture.CallFunction(t, "echo3", &s)
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

func TestHostFn_echo4_string(t *testing.T) {
	result, err := fixture.CallFunction(t, "echo4", "hello")
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

func TestHostFn_echo4_stringPtr(t *testing.T) {
	s := "hello"
	result, err := fixture.CallFunction(t, "echo4", &s)
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
