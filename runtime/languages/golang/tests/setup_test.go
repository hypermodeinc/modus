/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang_test

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/hypermodeinc/modus/runtime/testutils"
)

var basePath = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}()

var fixture *testutils.WasmTestFixture

func TestMain(m *testing.M) {
	path := filepath.Join(basePath, "..", "testdata", "build", "testdata.wasm")

	customTypes := make(map[string]reflect.Type)
	customTypes["testdata.TestStruct1"] = reflect.TypeFor[TestStruct1]()
	customTypes["testdata.TestStruct2"] = reflect.TypeFor[TestStruct2]()
	customTypes["testdata.TestStruct3"] = reflect.TypeFor[TestStruct3]()
	customTypes["testdata.TestStruct4"] = reflect.TypeFor[TestStruct4]()
	customTypes["testdata.TestStruct5"] = reflect.TypeFor[TestStruct5]()
	customTypes["testdata.TestStructWithMap"] = reflect.TypeFor[TestStructWithMap1]()
	customTypes["testdata.TestRecursiveStruct"] = reflect.TypeFor[TestRecursiveStruct]()

	registrations := getTestHostFunctionRegistrations()
	fixture = testutils.NewWasmTestFixture(path, customTypes, registrations)

	exitVal := m.Run()

	fixture.Close()
	os.Exit(exitVal)
}
