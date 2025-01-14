/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hypermodeinc/modus/runtime/hostfunctions"
	"github.com/hypermodeinc/modus/runtime/testutils"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
)

func getTestHostFunctionRegistrations() []func(wasmhost.WasmHost) error {
	return []func(wasmhost.WasmHost) error{
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("modus_system", "logMessage", hostLog)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("modus_test", "add", hostAdd)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("modus_test", "echo", hostEcho)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("modus_test", "echoObject", hostEchoObject)
		},
	}
}

func hostLog(ctx context.Context, level, message string) {
	if utils.DebugModeEnabled() {
		hostfunctions.LogMessage(ctx, level, message)
	}
	t := testutils.GetTestT(ctx)
	t.Logf("[%s] %s", level, message)
}

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
	fnName := "add"
	result, err := fixture.CallFunction(t, fnName, 1, 2)
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
	fnName := "echo"
	result, err := fixture.CallFunction(t, fnName, "hello")
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
	fnName := "echoObject"
	o := &TestHostObject{
		A: 1,
		B: true,
		C: "hello",
	}

	result, err := fixture.CallFunction(t, fnName, o)
	if err != nil {
		t.Fatal(err)
	}

	expected := &TestHostObject{
		A: 2,
		B: false,
		C: "hello!",
	}

	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(TestHostObject); !ok {
		t.Errorf("expected %T, got %T", expected, result)
	} else if reflect.DeepEqual(expected, r) {
		t.Errorf("expected %+v, got %+v", expected, r)
	}
}
