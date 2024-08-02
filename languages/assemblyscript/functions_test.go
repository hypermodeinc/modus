/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"hmruntime/plugins/metadata"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/wazero/api"
)

type WazeroOnly interface {
	wazeroOnly()
}

type WazeroOnlyType struct{}

func (WazeroOnlyType) wazeroOnly() {}

type MockModule struct {
	api.Module
	WazeroOnly
}

func (m *MockModule) String() string                                                 { return "" }
func (m *MockModule) Name() string                                                   { return "" }
func (m *MockModule) Memory() api.Memory                                             { return nil }
func (m *MockModule) ExportedFunction(name string) api.Function                      { return nil }
func (m *MockModule) ExportedFunctionDefinitions() map[string]api.FunctionDefinition { return nil }
func (m *MockModule) ExportedMemory(name string) api.Memory                          { return nil }
func (m *MockModule) ExportedMemoryDefinitions() map[string]api.MemoryDefinition     { return nil }
func (m *MockModule) ExportedGlobal(name string) api.Global                          { return nil }
func (m *MockModule) CloseWithExitCode(ctx context.Context, exitCode uint32) error   { return nil }
func (m *MockModule) IsClosed() bool                                                 { return false }
func (m *MockModule) Close(ctx context.Context) error                                { return m.CloseWithExitCode(ctx, 0) }

func Test_GetParameters(t *testing.T) {

	adapter := &wasmAdapter{}

	makeDefault := func(val any) *any {
		return &val
	}

	paramInfo := []*metadata.Parameter{
		{Name: "x", Type: "i32", Default: makeDefault(0)},
		{Name: "y", Type: "i32", Default: makeDefault(1)},
		{Name: "z", Type: "i32", Default: makeDefault(2)},
	}

	// no parameters supplied
	parameters := make(map[string]any)
	mockModule := &MockModule{}
	params, _, err := adapter.getParameters(context.Background(), mockModule, paramInfo, parameters)
	require.NoError(t, err)
	require.Equal(t, uint64(0), params[0])
	require.Equal(t, uint64(1), params[1])
	require.Equal(t, uint64(2), params[2])

	// only first parameter supplied
	parameters = make(map[string]any)
	parameters["x"] = 100
	mockModule = &MockModule{}
	params, _, err = adapter.getParameters(context.Background(), mockModule, paramInfo, parameters)
	require.NoError(t, err)
	require.Equal(t, uint64(100), params[0])
	require.Equal(t, uint64(1), params[1])
	require.Equal(t, uint64(2), params[2])

	// only second parameter supplied
	parameters = make(map[string]any)
	parameters["y"] = 100
	mockModule = &MockModule{}
	params, _, err = adapter.getParameters(context.Background(), mockModule, paramInfo, parameters)
	require.NoError(t, err)
	require.Equal(t, uint64(0), params[0])
	require.Equal(t, uint64(100), params[1])
	require.Equal(t, uint64(2), params[2])
}
