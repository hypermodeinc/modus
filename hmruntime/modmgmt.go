/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

var runtime wazero.Runtime

type buffers struct {
	Stdout *strings.Builder
	Stderr *strings.Builder
}

func initWasmRuntime(ctx context.Context) (wazero.Runtime, error) {

	// Create the runtime
	cfg := wazero.NewRuntimeConfig().
		WithCloseOnContextDone(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, cfg)

	// Connect WASI host functions
	err := instantiateWasiFunctions(ctx, runtime)
	if err != nil {
		return nil, err
	}

	// Connect Hypermode host functions
	err = instantiateHostFunctions(ctx, runtime)
	if err != nil {
		return nil, err
	}

	return runtime, nil
}

func loadPluginModule(ctx context.Context, name string) error {
	_, reloading := compiledModules[name]
	if reloading {
		fmt.Printf("Reloading plugin '%s'\n", name)
	} else {
		fmt.Printf("Loading plugin '%s'\n", name)
	}

	path, err := getPathForPlugin(name)
	if err != nil {
		return fmt.Errorf("failed to get path for plugin: %w", err)
	}

	plugin, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load the plugin: %w", err)
	}

	// Compile the plugin into a module.
	cm, err := runtime.CompileModule(ctx, plugin)
	if err != nil {
		return fmt.Errorf("failed to compile the plugin: %w", err)
	}

	// Store the compiled module for later retrieval.
	compiledModules[name] = cm

	// TODO: We should close the old module, but that leaves the _new_ module in an invalid state,
	// giving an error when querying: "source module must be compiled before instantiation"
	// if reloading {
	// 	cmOld.Close(ctx)
	// }

	return nil
}

func unloadPluginModule(ctx context.Context, name string) error {
	cmOld, found := compiledModules[name]
	if !found {
		return fmt.Errorf("plugin not found '%s'", name)
	}

	fmt.Printf("Unloading plugin '%s'\n", name)
	delete(compiledModules, name)
	cmOld.Close(ctx)

	return nil
}

func getModuleInstance(ctx context.Context, pluginName string) (wasm.Module, buffers, error) {

	// Create string buffers to capture stdout and stderr.
	// Still write to the console, but also capture the output in the buffers.
	buf := buffers{&strings.Builder{}, &strings.Builder{}}
	wOut := io.MultiWriter(os.Stdout, buf.Stdout)
	wErr := io.MultiWriter(os.Stderr, buf.Stderr)

	// Get the compiled module.
	compiled, ok := compiledModules[pluginName]
	if !ok {
		return nil, buf, fmt.Errorf("no compiled module found for plugin '%s'", pluginName)
	}

	// Configure the module instance.
	cfg := wazero.NewModuleConfig().
		WithName(pluginName + "_" + uuid.NewString()).
		WithStdout(wOut).WithStderr(wErr)

	// Instantiate the plugin as a module.
	// NOTE: This will also invoke the plugin's `_start` function,
	// which will call any top-level code in the plugin.
	mod, err := runtime.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		return nil, buf, fmt.Errorf("failed to instantiate the plugin module: %w", err)
	}

	return mod, buf, nil
}
