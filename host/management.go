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
	"strings"
	"time"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/logger"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type buffers struct {
	Stdout *strings.Builder
	Stderr *strings.Builder
}

func InitWasmRuntime(ctx context.Context) (wazero.Runtime, error) {

	// Create the runtime
	cfg := wazero.NewRuntimeConfig()
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

func loadPlugin(ctx context.Context, path string) (Plugin, error) {

	// Load the binary content of the plugin.
	bytes, err := getPluginBytes(ctx, path)
	if err != nil {
		return Plugin{}, err
	}

	// Compile the plugin into a module.
	cm, err := WasmRuntime.CompileModule(ctx, bytes)
	if err != nil {
		return Plugin{}, fmt.Errorf("failed to compile the plugin: %w", err)
	}

	// Get the metadata for the plugin.
	metadata, foundMetadata := getPluginMetadata(&cm)

	// Use the filename as the plugin name if no metadata is found.
	if !foundMetadata {
		metadata.Name, _ = getPluginNameFromPath(path)
	}

	// Create and store the plugin.
	plugin := Plugin{&cm, metadata, path}
	Plugins.Add(plugin)

	// Log the details of the loaded plugin.
	logPluginLoaded(ctx, plugin)
	if !foundMetadata {
		logger.Warn(ctx).
			Str("path", path).
			Str("plugin", plugin.Name()).
			Msg("No metadata found in plugin.  Please recompile your plugin using the latest version of the Hypermode Functions library.")
	}

	return plugin, nil
}

func logPluginLoaded(ctx context.Context, plugin Plugin) {
	evt := logger.Info(ctx)
	evt.Str("path", plugin.FilePath)

	metadata := plugin.Metadata

	if metadata.Name != "" {
		evt.Str("plugin", metadata.Name)
	}

	if metadata.Version != "" {
		evt.Str("version", metadata.Version)
	}

	lang := plugin.Language()
	if lang != UnknownLanguage {
		evt.Str("language", lang.String())
	}

	if metadata.BuildId != "" {
		evt.Str("build_id", metadata.BuildId)
	}

	if metadata.BuildTime != (time.Time{}) {
		evt.Time("build_id", metadata.BuildTime)
	}

	if metadata.LibraryName != "" {
		evt.Str("hypermode_library", metadata.LibraryName)
	}

	if metadata.LibraryVersion != "" {
		evt.Str("hypermode_library_version", metadata.LibraryVersion)
	}

	if metadata.GitRepo != "" {
		evt.Str("git_repo", metadata.GitRepo)
	}

	if metadata.GitCommit != "" {
		evt.Str("git_commit", metadata.GitCommit)
	}

	evt.Msg("Loaded plugin.")
}

func getPluginBytes(ctx context.Context, path string) ([]byte, error) {

	if aws.UseAwsForPluginStorage() {
		return aws.GetPluginBytes(ctx, path)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load the plugin: %w", err)
	}

	return bytes, nil
}

func unloadPlugin(ctx context.Context, plugin Plugin) {
	logger.Info(ctx).
		Str("plugin", plugin.Name()).
		Msg("Unloading plugin.")

	Plugins.Remove(plugin)
	(*plugin.Module).Close(ctx)
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
	plugin, ok := Plugins.GetByName(pluginName)
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
