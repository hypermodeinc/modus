/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

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

func NewGoWasmTestFixture(t *testing.T) *testutils.WasmTestFixture {
	path := filepath.Join(basePath, "..", "testdata", "testdata.wasm")
	registrations := getTestHostFunctionRegistrations()
	return testutils.NewWasmTestFixture(path, t, registrations...)
}
