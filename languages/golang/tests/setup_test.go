/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

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

	fixture.AddCustomType("testdata.TestStruct1", reflect.TypeFor[TestStruct1]())
	fixture.AddCustomType("testdata.TestStruct2", reflect.TypeFor[TestStruct2]())
	fixture.AddCustomType("testdata.TestStruct3", reflect.TypeFor[TestStruct3]())
	fixture.AddCustomType("testdata.TestStruct4", reflect.TypeFor[TestStruct4]())
	fixture.AddCustomType("testdata.TestStruct5", reflect.TypeFor[TestStruct5]())
	fixture.AddCustomType("testdata.TestStructWithMap", reflect.TypeFor[testStructWithMap]())
	fixture.AddCustomType("testdata.TestRecursiveStruct", reflect.TypeFor[TestRecursiveStruct]())

	exitVal := m.Run()

	fixture.Close()
	os.Exit(exitVal)
}
