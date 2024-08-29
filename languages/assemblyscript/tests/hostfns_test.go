/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"context"
	"reflect"
	"testing"

	"hmruntime/testutils"
	"hmruntime/wasmhost"
)

func getTestHostFunctionRegistrations() []func(*wasmhost.WasmHost) error {
	return []func(*wasmhost.WasmHost) error{
		func(host *wasmhost.WasmHost) error {
			return host.RegisterHostFunction("hypermode", "log", hostLog)
		},
		func(host *wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "add", hostAdd)
		},
		func(host *wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "echo", hostEcho)
		},
		func(host *wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "echoObject", hostEchoObject)
		},
	}
}

func hostLog(ctx context.Context, level, message string) {
	t := testutils.GetTestT(ctx)
	t.Logf("[%s] %s", level, message)
}

// TODO: we should be able to pass these as int
func hostAdd(a, b int32) int32 {
	return a + b
}

func hostEcho(s string) string {
	return "echo: " + s
}

func hostEchoObject(obj *TestHostObject) *TestHostObject {
	return &TestHostObject{
		A: obj.A + 1,
		B: !obj.B,
		C: obj.C + "!",
	}
}

type TestHostObject struct {
	A int32
	B bool
	C string
}

func TestHostFn_add(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	result, err := f.InvokeFunction("add", 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(int32); !ok {
		t.Errorf("expected an int32, got %T", result)
	} else if r != 3 {
		t.Errorf("expected %d, got %d", 3, r)
	}
}

func TestHostFn_echo(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
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

func TestHostFn_echoObject(t *testing.T) {
	t.Parallel()

	f := NewASWasmTestFixture(t)
	defer f.Close()

	f.AddCustomType("assembly/hostfns/TestHostObject", reflect.TypeFor[TestHostObject]())

	o := &TestHostObject{
		A: 1,
		B: true,
		C: "hello",
	}

	result, err := f.InvokeFunction("echoObject", o)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]any{
		"A": int32(2),
		"B": false,
		"C": "hello!",
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(map[string]any); !ok {
		t.Errorf("expected a map[string]any, got %T", result)
	} else if reflect.DeepEqual(expected, r) {
		t.Errorf("expected %+v, got %+v", expected, r)
	}
}
