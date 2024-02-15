/*
 * Copyright 2023 Hypermode, Inc.
 */

package plugins

import (
	"context"
	"fmt"
	"hmruntime/aws"
	"hmruntime/functions"
	"hmruntime/host"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental/opt"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type buffers struct {
	Stdout *strings.Builder
	Stderr *strings.Builder
}

func InitWasmRuntime(ctx context.Context) (wazero.Runtime, error) {

	// Create the runtime
	var cfg wazero.RuntimeConfig
	if runtime.GOARCH == "arm64" {
		// Use the experimental optimizing compiler for ARM64 to improve performance.
		// This is not yet available for other architectures.
		// See https://github.com/tetratelabs/wazero/releases/tag/v1.6.0
		cfg = opt.NewRuntimeConfigOptimizingCompiler()
	} else {
		cfg = wazero.NewRuntimeConfig()
	}

	cfg = cfg.WithCloseOnContextDone(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, cfg)

	// Connect WASI host functions
	err := instantiateWasiFunctions(ctx, runtime)
	if err != nil {
		return nil, err
	}

	// Connect Hypermode host functions
	err = functions.InstantiateHostFunctions(ctx, runtime)
	if err != nil {
		return nil, err
	}

	return runtime, nil
}

func loadPluginModule(ctx context.Context, name string) error {
	_, reloading := host.CompiledModules[name]
	log.Info().
		Str("plugin", name).
		Bool("reloading", reloading).
		Msg("Loading plugin.")

	// Load the binary content of the plugin.
	plugin, err := getPluginBytes(ctx, name)
	if err != nil {
		return err
	}

	// Compile the plugin into a module.
	cm, err := host.WasmRuntime.CompileModule(ctx, plugin)
	if err != nil {
		return fmt.Errorf("failed to compile the plugin: %w", err)
	}

	// Store the compiled module for later retrieval.
	host.CompiledModules[name] = cm

	// TODO: We should close the old module, but that leaves the _new_ module in an invalid state,
	// giving an error when querying: "source module must be compiled before instantiation"
	// if reloading {
	// 	cmOld.Close(ctx)
	// }

	return nil
}

func getPluginBytes(ctx context.Context, name string) ([]byte, error) {

	if aws.UseAwsForPluginStorage() {
		return aws.GetPluginBytes(ctx, name)
	}

	path, err := getPathForPlugin(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get path for plugin: %w", err)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load the plugin: %w", err)
	}

	log.Info().
		Str("plugin", name).
		Str("path", path).
		Msg("Retrieved plugin from local storage.")

	return bytes, nil
}

func unloadPluginModule(ctx context.Context, name string) error {
	cmOld, found := host.CompiledModules[name]
	if !found {
		return fmt.Errorf("plugin not found '%s'", name)
	}

	log.Info().
		Str("plugin", name).
		Msg("Unloading plugin.")

	delete(host.CompiledModules, name)
	cmOld.Close(ctx)

	return nil
}

func GetModuleInstance(ctx context.Context, pluginName string) (wasm.Module, buffers, error) {

	// Create string buffers to capture stdout and stderr.
	// Still write to the console, but also capture the output in the buffers.
	buf := buffers{&strings.Builder{}, &strings.Builder{}}
	wOut := io.MultiWriter(os.Stdout, buf.Stdout)
	wErr := io.MultiWriter(os.Stderr, buf.Stderr)

	// Get the compiled module.
	compiled, ok := host.CompiledModules[pluginName]
	if !ok {
		return nil, buf, fmt.Errorf("no compiled module found for plugin '%s'", pluginName)
	}

	// Configure the module instance.
	cfg := wazero.NewModuleConfig().
		WithName(pluginName + "_" + uuid.NewString()).
		WithSysWalltime().WithSysNanotime().
		WithStdout(wOut).WithStderr(wErr)

	// Instantiate the plugin as a module.
	// NOTE: This will also invoke the plugin's `_start` function,
	// which will call any top-level code in the plugin.
	mod, err := host.WasmRuntime.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		return nil, buf, fmt.Errorf("failed to instantiate the plugin module: %w", err)
	}

	return mod, buf, nil
}

func instantiateWasiFunctions(ctx context.Context, runtime wazero.Runtime) error {
	b := runtime.NewHostModuleBuilder(wasi.ModuleName)
	wasi.NewFunctionExporter().ExportFunctions(b)

	// If we ever need to override any of the WASI functions, we can do so here.

	_, err := b.Instantiate(ctx)
	if err != nil {
		return fmt.Errorf("failed to instantiate the %s module: %w", wasi.ModuleName, err)
	}

	return nil
}
