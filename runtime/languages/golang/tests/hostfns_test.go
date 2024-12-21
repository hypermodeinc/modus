/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang_test

import (
	"context"
	"strings"
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
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "encodeStrings1", hostEncodeStrings1)
		},
		func(host wasmhost.WasmHost) error {
			return host.RegisterHostFunction("test", "encodeStrings2", hostEncodeStrings2)
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

func hostEncodeStrings1(items *[]string) *string {
	out := strings.Builder{}
	out.WriteString("[")
	for i, item := range *items {
		if i > 0 {
			out.WriteString(",")
		}
		b, _ := utils.JsonSerialize(item)
		out.Write(b)
	}
	out.WriteString("]")

	result := out.String()
	return &result
}

func hostEncodeStrings2(items *[]*string) *string {
	out := strings.Builder{}
	out.WriteString("[")
	for i, item := range *items {
		if i > 0 {
			out.WriteString(",")
		}
		b, _ := utils.JsonSerialize(item)
		out.Write(b)
	}
	out.WriteString("]")

	result := out.String()
	return &result
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

func TestHostFn_echo1_string(t *testing.T) {
	fnName := "echo1"
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

func TestHostFn_echo1_stringPtr(t *testing.T) {
	fnName := "echo1"
	s := "hello"

	result, err := fixture.CallFunction(t, fnName, &s)
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
	fnName := "echo2"
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

func TestHostFn_echo2_stringPtr(t *testing.T) {
	fnName := "echo2"
	s := "hello"

	result, err := fixture.CallFunction(t, fnName, &s)
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
	fnName := "echo3"
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

func TestHostFn_echo3_stringPtr(t *testing.T) {
	fnName := "echo3"
	s := "hello"

	result, err := fixture.CallFunction(t, fnName, &s)
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
	fnName := "echo4"
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

func TestHostFn_echo4_stringPtr(t *testing.T) {
	fnName := "echo4"
	s := "hello"

	result, err := fixture.CallFunction(t, fnName, &s)
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

func TestHostFn_encodeStrings1(t *testing.T) {
	fnName := "encodeStrings1"
	s := []string{"hello", "world"}

	result, err := fixture.CallFunction(t, fnName, s)
	if err != nil {
		t.Fatal(err)
	}

	expected := `["hello","world"]`
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(string); !ok {
		t.Errorf("expected a string, got %T", result)
	} else if r != expected {
		t.Errorf("expected %s, got %s", expected, r)
	}
}

func TestHostFn_encodeStrings2(t *testing.T) {
	fnName := "encodeStrings2"
	e0 := "hello"
	e1 := "world"
	s := []*string{&e0, &e1}

	result, err := fixture.CallFunction(t, fnName, s)
	if err != nil {
		t.Fatal(err)
	}

	expected := `["hello","world"]`
	if result == nil {
		t.Error("expected a result")
	} else if r, ok := result.(string); !ok {
		t.Errorf("expected a string, got %T", result)
	} else if r != expected {
		t.Errorf("expected %s, got %s", expected, r)
	}
}
