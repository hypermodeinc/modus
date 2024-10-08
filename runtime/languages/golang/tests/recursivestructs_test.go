/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang_test

type TestRecursiveStruct struct {
	A bool
	B *TestRecursiveStruct
}

// testRecursiveStruct := func() TestRecursiveStruct {
// 	r := TestRecursiveStruct{
// 		A: true,
// 	}
// 	r.B = &r
// 	return r
// }()

// testRecursiveStructAsMap := func() map[string]any {
// 	r1 := map[string]any{
// 		"a": true,
// 	}
// 	r2 := map[string]any{
// 		"a": false,
// 	}
// 	r1["b"] = r2
// 	r2["b"] = r1
// 	return r1
// }()

// func TestRecursiveStructInput(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	if _, err := f.InvokeFunction("testRecursiveStructInput", testRecursiveStruct); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestRecursiveStructPtrInput(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	if _, err := f.InvokeFunction("testRecursiveStructPtrInput", testRecursiveStruct); err != nil {
// 		t.Fatal(err)
// 	}
// 	if _, err := f.InvokeFunction("testRecursiveStructPtrInput", &testRecursiveStruct); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestRecursiveStructInput_map(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	if _, err := f.InvokeFunction("testRecursiveStructInput", testRecursiveStructAsMap); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestRecursiveStructPtrInput_map(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	if _, err := f.InvokeFunction("testRecursiveStructInput", testRecursiveStructAsMap); err != nil {
// 		t.Fatal(err)
// 	}
// 	if _, err := f.InvokeFunction("testRecursiveStructInput", &testRecursiveStructAsMap); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestRecursiveStructPtrInput_nil(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	if _, err := f.InvokeFunction("testRecursiveStructPtrInput_nil", nil); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestRecursiveStructOutput(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	fixture.AddCustomType("testdata.TestRecursiveStruct", reflect.TypeFor[TestRecursiveStruct]())

// 	result, err := f.InvokeFunction("testRecursiveStructOutput")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result == nil {
// 		t.Error("expected a result")
// 	} else if r, ok := result.(TestRecursiveStruct); !ok {
// 		t.Errorf("expected a struct, got %T", result)
// 	} else if !reflect.DeepEqual(testRecursiveStruct, r) {
// 		t.Errorf("expected %v, got %v", testRecursiveStruct, r)
// 	}
// }

// func TestRecursiveStructPtrOutput(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	fixture.AddCustomType("testdata.TestRecursiveStruct", reflect.TypeFor[TestRecursiveStruct]())

// 	result, err := f.InvokeFunction("testRecursiveStructPtrOutput")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result == nil {
// 		t.Error("expected a result")
// 	} else if r, ok := result.(*TestRecursiveStruct); !ok {
// 		t.Errorf("expected a pointer to a struct, got %T", result)
// 	} else if !reflect.DeepEqual(testRecursiveStruct, *r) {
// 		t.Errorf("expected %v, got %v", testRecursiveStruct, *r)
// 	}
// }

// func TestRecursiveStructOutput_map(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	result, err := f.InvokeFunction("testRecursiveStructOutput")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result == nil {
// 		t.Error("expected a result")
// 	} else if r, ok := result.(map[string]any); !ok {
// 		t.Errorf("expected a map[string]any, got %T", result)
// 	} else if !reflect.DeepEqual(testRecursiveStructAsMap, r) {
// 		t.Errorf("expected %v, got %v", testRecursiveStructAsMap, r)
// 	}
// }

// func TestRecursiveStructPtrOutput_map(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	result, err := f.InvokeFunction("testRecursiveStructPtrOutput")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result == nil {
// 		t.Error("expected a result")
// 	} else if r, ok := result.(map[string]any); !ok {
// 		t.Errorf("expected a map[string]any, got %T", result)
// 	} else if !reflect.DeepEqual(testRecursiveStructAsMap, r) {
// 		t.Errorf("expected %v, got %v", testRecursiveStructAsMap, r)
// 	}
// }

// func TestRecursiveStructPtrOutput_nil(t *testing.T) {
// // 	f := NewGoWasmTestFixture(t)
// 	defer f.Close()

// 	result, err := f.InvokeFunction("testRecursiveStructPtrOutput_nil")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !utils.HasNil(result) {
// 		t.Error("expected a nil result")
// 	}
// }
