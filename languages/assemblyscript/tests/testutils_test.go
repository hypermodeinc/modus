/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"hypruntime/testutils"
)

var basePath = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}()

func NewASWasmTestFixture(t *testing.T) *testutils.WasmTestFixture {
	path := filepath.Join(basePath, "..", "testdata", "build", "testdata.wasm")
	registrations := getTestHostFunctionRegistrations()
	return testutils.NewWasmTestFixture(path, t, registrations...)
}
