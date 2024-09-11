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
	registrations := getTestHostFunctionRegistrations()
	fixture = testutils.NewWasmTestFixture(path, registrations...)

	fixture.AddCustomType("assembly/classes/TestClass1", reflect.TypeFor[TestClass1]())
	fixture.AddCustomType("assembly/classes/TestClass2", reflect.TypeFor[TestClass2]())
	fixture.AddCustomType("assembly/classes/TestClass3", reflect.TypeFor[TestClass3]())
	fixture.AddCustomType("assembly/classes/TestClass4", reflect.TypeFor[TestClass4]())
	fixture.AddCustomType("assembly/classes/TestClass5", reflect.TypeFor[TestClass5]())
	fixture.AddCustomType("assembly/hostfns/TestHostObject", reflect.TypeFor[TestHostObject]())
	fixture.AddCustomType("assembly/maps/TestClassWithMap", reflect.TypeFor[testClassWithMap]())

	exitVal := m.Run()

	fixture.Close()
	os.Exit(exitVal)
}
