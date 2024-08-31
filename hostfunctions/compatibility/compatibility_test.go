/*
 * Copyright 2024 Hypermode, Inc.
 */

package compatibility_test

import (
	"context"
	"testing"

	"hypruntime/functions"
	"hypruntime/hostfunctions"
	"hypruntime/hostfunctions/compatibility"
	"hypruntime/plugins"
	"hypruntime/utils"
	"hypruntime/wasmhost"

	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

func Test_HostFunctionsHaveCompatibilityShim(t *testing.T) {

	registrations := hostfunctions.GetRegistrations()

	host := &mockWasmHost{}
	for _, reg := range registrations {
		if err := reg(host); err != nil {
			t.Errorf("Error registering host function: %v", err)
		}
	}

	for _, name := range host.hostFunctions {
		fn, err := compatibility.GetImportMetadataShim(name)
		if err != nil {
			t.Errorf("Missing import metadata shim for host function: %v", err)
		} else if fn == nil {
			t.Errorf("Invalid import metadata shim for host function: %v", name)
		}
	}

}

type mockWasmHost struct {
	hostFunctions []string
}

func (m *mockWasmHost) RegisterHostFunction(modName, funcName string, fn any, opts ...wasmhost.HostFunctionOption) error {
	m.hostFunctions = append(m.hostFunctions, modName+"."+funcName)
	return nil
}

func (m *mockWasmHost) CallFunction(ctx context.Context, fnInfo functions.FunctionInfo, parameters map[string]any) (wasmhost.ExecutionInfo, error) {
	panic("not implemented")
}
func (m *mockWasmHost) CallFunctionByName(ctx context.Context, fnName string, paramValues ...any) (wasmhost.ExecutionInfo, error) {
	panic("not implemented")
}
func (m *mockWasmHost) Close(ctx context.Context) {
	panic("not implemented")
}
func (m *mockWasmHost) CompileModule(ctx context.Context, bytes []byte) (wazero.CompiledModule, error) {
	panic("not implemented")
}
func (m *mockWasmHost) GetFunctionInfo(fnName string) (functions.FunctionInfo, error) {
	panic("not implemented")
}
func (m *mockWasmHost) GetFunctionRegistry() functions.FunctionRegistry {
	panic("not implemented")
}
func (m *mockWasmHost) GetModuleInstance(ctx context.Context, plugin *plugins.Plugin, buffers utils.OutputBuffers) (wasm.Module, error) {
	panic("not implemented")
}
