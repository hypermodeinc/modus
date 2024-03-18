/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"hmruntime/config"
	"hmruntime/logger"
	"hmruntime/storage"

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

func InitWasmRuntime(ctx context.Context) error {
	cfg := wazero.NewRuntimeConfig()
	cfg = cfg.WithCloseOnContextDone(true)
	WasmRuntime = wazero.NewRuntimeWithConfig(ctx, cfg)
	return instantiateWasiFunctions(ctx)
}

func loadJson(ctx context.Context, filename string) error {
	bytes, err := storage.GetFileContents(ctx, filename)
	if err != nil {
		return err
	}

	_, ok := config.SupportedJsons[filename]
	if ok {
		err = json.Unmarshal(bytes, config.SupportedJsons[filename])
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", filename, err)
		}
	}

	return nil
}

func loadPlugin(ctx context.Context, filename string) (Plugin, error) {

	// Load the binary content of the plugin.
	bytes, err := storage.GetFileContents(ctx, filename)
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
		metadata.Name, _ = getPluginNameFromPath(filename)
	}

	// Create and store the plugin.
	plugin := Plugin{&cm, metadata, filename}
	Plugins.Add(plugin)

	// Log the details of the loaded plugin.
	logPluginLoaded(ctx, plugin)
	if !foundMetadata {
		logger.Warn(ctx).
			Str("filename", filename).
			Str("plugin", plugin.Name()).
			Msg("No metadata found in plugin.  Please recompile your plugin using the latest version of the Hypermode Functions library.")
	}

	return plugin, nil
}

func logPluginLoaded(ctx context.Context, plugin Plugin) {
	evt := logger.Info(ctx)
	evt.Str("filename", plugin.FileName)

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

func instantiateWasiFunctions(ctx context.Context) error {
	b := WasmRuntime.NewHostModuleBuilder(wasi.ModuleName)
	wasi.NewFunctionExporter().ExportFunctions(b)

	// If we ever need to override any of the WASI functions, we can do so here.

	_, err := b.Instantiate(ctx)
	if err != nil {
		return fmt.Errorf("failed to instantiate the %s module: %w", wasi.ModuleName, err)
	}

	return nil
}
