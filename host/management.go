/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/logger"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
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

	return runtime, nil
}

func loadJson(ctx context.Context, name string) error {
	config.Mu.Lock()         // Lock the mutex
	defer config.Mu.Unlock() // Unlock the mutex when the function returns
	logger.Info(ctx).Str("Loading %s.json.", name)

	// Get the JSON bytes
	bytes, err := getJsonBytes(ctx, name)
	if err != nil {
		return err
	}

	_, exists := config.SupportedJsons[name+".json"]
	if !exists {
		return fmt.Errorf("JSON %s does not exist.", name)
	}

	err = json.Unmarshal(bytes, config.SupportedJsons[name+".json"])
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s.json: %w", name, err)
	}

	return nil
}

func unloadJson(ctx context.Context, name string) {
	config.Mu.Lock()         // Lock the mutex
	defer config.Mu.Unlock() // Unlock the mutex when the function returns

	value, exists := config.SupportedJsons[name+".json"]
	if !exists {
		logger.Error(ctx).Msg(fmt.Sprintf("JSON %s does not exist.", name))
		return
	}

	t := reflect.TypeOf(value).Elem()
	v := reflect.New(t).Interface()
	config.SupportedJsons[name+".json"] = v

	logger.Info(ctx).Msg(fmt.Sprintf("Unloaded %s.json.", name))
}

func getJsonBytes(ctx context.Context, name string) ([]byte, error) {
	if aws.UseAwsForPluginStorage() {
		return aws.GetJsonBytes(ctx, name)
	}

	path, err := getPathForJson(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get path for %s.json: %w", name, err)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s.json: %w", name, err)
	}

	logger.Info(ctx).
		Str("path", path).
		Msg(fmt.Sprintf("Retrieved %s.json from local storage.", name))

	return bytes, nil
}

func loadPlugin(ctx context.Context, name string) error {
	_, reloading := Plugins[name]
	logger.Info(ctx).
		Str("filename", name+".wasm").
		Bool("reloading", reloading).
		Msg("Loading plugin.")

	// TODO: Separate plugin name from file name throughout the codebase.
	// The plugin name should always come from the metadata, not the file name.
	// This requires significant changes, so it's not done yet.

	// Load the binary content of the plugin.
	bytes, err := getPluginBytes(ctx, name)
	if err != nil {
		return err
	}

	// Compile the plugin into a module.
	cm, err := WasmRuntime.CompileModule(ctx, bytes)
	if err != nil {
		return fmt.Errorf("failed to compile the plugin: %w", err)
	}

	// Get the metadata for the plugin.
	metadata := getPluginMetadata(&cm)
	if metadata == (PluginMetadata{}) {
		logger.Warn(ctx).
			Str("filename", name+".wasm").
			Msg("No metadata found in plugin.  Please recompile your plugin using the latest version of the Hypermode Functions library.")
	}

	// Finally, store the plugin to complete the loading process.
	Plugins[name] = Plugin{&cm, metadata}
	logPluginLoaded(ctx, metadata)

	return nil
}

func logPluginLoaded(ctx context.Context, metadata PluginMetadata) {
	evt := logger.Info(ctx)

	if metadata != (PluginMetadata{}) {
		evt.Str("plugin", metadata.Name)
		evt.Str("version", metadata.Version)
		evt.Str("language", metadata.Language.String())
		evt.Str("build_id", metadata.BuildId)
		evt.Time("build_time", metadata.BuildTime)
		evt.Str("hypermode_library", metadata.LibraryName)
		evt.Str("hypermode_library_version", metadata.LibraryVersion)

		if metadata.GitRepo != "" {
			evt.Str("git_repo", metadata.GitRepo)
		}

		if metadata.GitCommit != "" {
			evt.Str("git_commit", metadata.GitCommit)
		}
	}

	evt.Msg("Loaded plugin.")
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

	logger.Info(ctx).
		Str("plugin", name).
		Str("path", path).
		Msg("Retrieved plugin from local storage.")

	return bytes, nil
}

func unloadPlugin(ctx context.Context, name string) error {
	plugin, found := Plugins[name]
	if !found {
		return fmt.Errorf("plugin not found '%s'", name)
	}

	logger.Info(ctx).
		Str("plugin", name).
		Msg("Unloading plugin.")

	mod := *plugin.Module
	delete(Plugins, name)
	mod.Close(ctx)

	return nil
}

func GetModuleInstance(ctx context.Context, pluginName string) (wasm.Module, buffers, error) {

	// Get the logger and writers for the plugin's stdout and stderr.
	log := logger.Get(ctx).With().Bool("user_visible", true).Logger()
	wInfoLog := logger.NewLogWriter(&log, zerolog.InfoLevel)
	wErrorLog := logger.NewLogWriter(&log, zerolog.ErrorLevel)

	// Create string buffers to capture stdout and stderr.
	// Still write to the log, but also capture the output in the buffers.
	buf := buffers{&strings.Builder{}, &strings.Builder{}}
	wOut := io.MultiWriter(buf.Stdout, wInfoLog)
	wErr := io.MultiWriter(buf.Stderr, wErrorLog)

	// Get the plugin.
	plugin, ok := Plugins[pluginName]
	if !ok {
		return nil, buf, fmt.Errorf("plugin not found with name '%s'", pluginName)
	}

	// Configure the module instance.
	cfg := wazero.NewModuleConfig().
		WithName(pluginName + "_" + uuid.NewString()).
		WithSysWalltime().WithSysNanotime().
		WithStdout(wOut).WithStderr(wErr)

	// Instantiate the plugin as a module.
	// NOTE: This will also invoke the plugin's `_start` function,
	// which will call any top-level code in the plugin.
	mod, err := WasmRuntime.InstantiateModule(ctx, *plugin.Module, cfg)
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
