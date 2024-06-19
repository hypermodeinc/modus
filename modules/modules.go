/*
 * Copyright 2024 Hypermode, Inc.
 */

package modules

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"

	"hmruntime/functions"
	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/utils"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

// Global runtime instance for the WASM modules
var RuntimeInstance wazero.Runtime

// Global, thread-safe registry of all plugins loaded by the host
var Plugins = plugins.NewPluginRegistry()

// Gets a module instance for the given plugin, used for a single invocation.
func GetModuleInstance(ctx context.Context, plugin *plugins.Plugin, buffers *utils.OutputBuffers) (wasm.Module, error) {

	// Get the logger and writers for the plugin's stdout and stderr.
	log := logger.Get(ctx).With().Bool("user_visible", true).Logger()
	wInfoLog := logger.NewLogWriter(&log, zerolog.InfoLevel)
	wErrorLog := logger.NewLogWriter(&log, zerolog.ErrorLevel)

	// Capture stdout/stderr both to logs, and to provided writers.
	wOut := io.MultiWriter(&buffers.StdOut, wInfoLog)
	wErr := io.MultiWriter(&buffers.StdErr, wErrorLog)

	// Configure the module instance.
	cfg := wazero.NewModuleConfig().
		WithName(plugin.Name() + "_" + uuid.NewString()).
		WithSysWalltime().WithSysNanotime().
		WithRandSource(rand.Reader).
		WithStdout(wOut).WithStderr(wErr)

	// Instantiate the plugin as a module.
	// NOTE: This will also invoke the plugin's `_start` function,
	// which will call any top-level code in the plugin.
	mod, err := RuntimeInstance.InstantiateModule(ctx, *plugin.Module, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the plugin module: %w", err)
	}

	return mod, nil
}

func Setup(ctx context.Context, info functions.FunctionInfo) (context.Context, string, utils.OutputBuffers, []utils.LogMessage) {
	// Prepare the context that will be used throughout the function execution
	executionId := xid.New().String()
	ctx = context.WithValue(ctx, utils.ExecutionIdContextKey, executionId)
	ctx = context.WithValue(ctx, utils.PluginContextKey, info.Plugin)

	// Also prepare a slice to capture log messages sent through the "log" host function.
	messages := []utils.LogMessage{}
	ctx = context.WithValue(ctx, utils.FunctionMessagesContextKey, &messages)

	// Create output buffers for the function to write stdout/stderr to
	buffers := utils.OutputBuffers{}

	return ctx, executionId, buffers, messages
}

func CallFunctionByName(ctx context.Context, fnName string, inputValues ...any) (any, error) {
	info, ok := functions.Functions[fnName]
	if !ok {
		return nil, fmt.Errorf("no function registered named %s", fnName)
	}

	ctx, _, buffers, _ := Setup(ctx, info)

	mod, err := GetModuleInstance(ctx, info.Plugin, &buffers)
	if err != nil {
		return nil, err
	}
	defer mod.Close(ctx)

	return CallFunctionByNameWithModule(ctx, mod, fnName, inputValues...)
}

func CallFunctionByNameWithModule(ctx context.Context, mod wasm.Module, fnName string, inputValues ...any) (any, error) {
	info, ok := functions.Functions[fnName]
	if !ok {
		return nil, fmt.Errorf("no function registered named %s", fnName)
	}

	parameters := make(map[string]any, len(inputValues))
	for i, value := range inputValues {
		name := info.Function.Parameters[i].Name
		parameters[name] = value
	}

	return functions.CallFunction(ctx, mod, info, parameters)
}

func VerifyFunctionSignature(fnName string, expectedTypes ...string) error {
	info, ok := functions.Functions[fnName]
	if !ok {
		return fmt.Errorf("no function registered named %s", fnName)
	}

	if len(expectedTypes) == 0 {
		return errors.New("expectedTypes must not be empty")
	}
	l := len(expectedTypes)
	expectedSig := fmt.Sprintf("(%s):%s", strings.Join(expectedTypes[:l-1], ","), expectedTypes[l-1])

	sig := info.Function.Signature()
	if sig != expectedSig {
		return fmt.Errorf("function %s has signature %s, expected %s", fnName, sig, expectedSig)
	}

	return nil
}
