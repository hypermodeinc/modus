/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"hypruntime/testutils"
)

var basePath = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}()

var fixture *testutils.WasmTestFixture

func TestMain(m *testing.M) {
	path := filepath.Join(basePath, "..", "testdata", "build", "testdata.wasm")

	customTypes := make(map[string]reflect.Type)
	customTypes["assembly/classes/TestClass1"] = reflect.TypeFor[TestClass1]()
	customTypes["assembly/classes/TestClass2"] = reflect.TypeFor[TestClass2]()
	customTypes["assembly/classes/TestClass3"] = reflect.TypeFor[TestClass3]()
	customTypes["assembly/classes/TestClass4"] = reflect.TypeFor[TestClass4]()
	customTypes["assembly/classes/TestClass5"] = reflect.TypeFor[TestClass5]()
	customTypes["assembly/hostfns/TestHostObject"] = reflect.TypeFor[TestHostObject]()
	customTypes["assembly/maps/TestClassWithMap"] = reflect.TypeFor[TestClassWithMap1]()

	registrations := getTestHostFunctionRegistrations()
	fixture = testutils.NewWasmTestFixture(path, customTypes, registrations)

	exitVal := m.Run()

	fixture.Close()
	os.Exit(exitVal)
}
