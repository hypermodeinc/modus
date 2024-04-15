/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/storage"
	"hmruntime/utils"

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

var pluginLoaded func(ctx context.Context, metadata plugins.PluginMetadata) error

func RegisterPluginLoadedCallback(callback func(ctx context.Context, metadata plugins.PluginMetadata) error) {
	pluginLoaded = callback
}

func InitWasmRuntime(ctx context.Context) error {
	cfg := wazero.NewRuntimeConfig()
	cfg = cfg.WithCloseOnContextDone(true)
	WasmRuntime = wazero.NewRuntimeWithConfig(ctx, cfg)
	return instantiateWasiFunctions(ctx)
}

func MonitorPlugins(ctx context.Context) {

	loadPluginFile := func(fi storage.FileInfo) {
		err := loadPlugin(ctx, fi.Name)
		if err != nil {
			logger.Err(ctx, err).
				Str("filename", fi.Name).
				Msg("Failed to load plugin.")
		}
	}

	sm := storage.NewStorageMonitor(".wasm")
	sm.Added = loadPluginFile
	sm.Modified = loadPluginFile
	sm.Removed = func(fi storage.FileInfo) {
		p, ok := Plugins.GetByFile(fi.Name)
		if ok {
			err := unloadPlugin(ctx, p)
			if err != nil {
				logger.Err(ctx, err).
					Str("filename", fi.Name).
					Msg("Failed to unload plugin.")
			}
		}
	}
	sm.Changed = func() {
		// Signal that we need to register functions
		RegistrationRequest <- true
	}
	sm.Start(ctx)
}

func loadPlugin(ctx context.Context, filename string) error {

	// Load the binary content of the plugin.
	bytes, err := storage.GetFileContents(ctx, filename)
	if err != nil {
		return err
	}

	// Compile the plugin into a module.
	cm, err := WasmRuntime.CompileModule(ctx, bytes)
	if err != nil {
		return fmt.Errorf("failed to compile the plugin: %w", err)
	}

	// Get the metadata for the plugin.
	metadata, err := plugins.GetPluginMetadata(&cm)

	// Apply fallback rules for missing metadata.
	if err == plugins.ErrPluginMetadataNotFound {
		var found bool
		metadata, found = plugins.GetPluginMetadata_old(&cm)
		if found {
			defer logger.Warn(ctx).
				Str("filename", filename).
				Str("plugin", metadata.Plugin).
				Msg("Deprecated metadata format found in plugin.  Please recompile your plugin using the latest version of the Hypermode Functions library.")
		} else {
			// Use the filename as the plugin name if no metadata is found.
			metadata.Plugin = strings.TrimSuffix(filename, ".wasm")
			defer logger.Warn(ctx).
				Str("filename", filename).
				Str("plugin", metadata.Plugin).
				Msg("No metadata found in plugin.  Please recompile your plugin using the latest version of the Hypermode Functions library.")
		}
	} else if err != nil {
		return err
	}

	// Store the types in a map for easy access.
	types := make(map[string]plugins.TypeDefinition, len(metadata.Types))
	for _, t := range metadata.Types {
		types[t.Path] = t
	}

	// Create and store the plugin.
	plugin := plugins.Plugin{Module: &cm, Metadata: metadata, FileName: filename, Types: types}
	Plugins.Add(plugin)

	// Log the details of the loaded plugin.
	logPluginLoaded(ctx, plugin)

	// Notify the callback that a plugin has been loaded.
	err = pluginLoaded(ctx, metadata)
	if err != nil {
		return fmt.Errorf("plugin loaded callback failed: %w", err)
	}

	return nil
}

func logPluginLoaded(ctx context.Context, plugin plugins.Plugin) {
	evt := logger.Info(ctx)
	evt.Str("filename", plugin.FileName)

	metadata := plugin.Metadata

	if metadata.Plugin != "" {
		name, version := utils.ParseNameAndVersion(metadata.Plugin)
		evt.Str("plugin", name)
		evt.Str("version", version)
	}

	lang := plugin.Language()
	if lang != plugins.UnknownLanguage {
		evt.Str("language", lang.String())
	}

	if metadata.BuildId != "" {
		evt.Str("build_id", metadata.BuildId)
	}

	if metadata.BuildTime != (time.Time{}) {
		evt.Time("build_ts", metadata.BuildTime)
	}

	if metadata.Library != "" {
		name, version := utils.ParseNameAndVersion(metadata.Library)
		evt.Str("hypermode_library", name)
		evt.Str("hypermode_library_version", version)
	}

	if metadata.GitRepo != "" {
		evt.Str("git_repo", metadata.GitRepo)
	}

	if metadata.GitCommit != "" {
		evt.Str("git_commit", metadata.GitCommit)
	}

	evt.Msg("Loaded plugin.")
}

func unloadPlugin(ctx context.Context, plugin plugins.Plugin) error {
	logger.Info(ctx).
		Str("plugin", plugin.Name()).
		Str("build_id", plugin.BuildId()).
		Msg("Unloading plugin.")

	Plugins.Remove(plugin)
	return (*plugin.Module).Close(ctx)
}

func GetModuleInstance(ctx context.Context, plugin *plugins.Plugin) (wasm.Module, buffers, error) {

	// Get the logger and writers for the plugin's stdout and stderr.
	log := logger.Get(ctx).With().Bool("user_visible", true).Logger()
	wInfoLog := logger.NewLogWriter(&log, zerolog.InfoLevel)
	wErrorLog := logger.NewLogWriter(&log, zerolog.ErrorLevel)

	// Create string buffers to capture stdout and stderr.
	// Still write to the log, but also capture the output in the buffers.
	buf := buffers{&strings.Builder{}, &strings.Builder{}}
	wOut := io.MultiWriter(buf.Stdout, wInfoLog)
	wErr := io.MultiWriter(buf.Stderr, wErrorLog)

	// Configure the module instance.
	cfg := wazero.NewModuleConfig().
		WithName(plugin.Name() + "_" + uuid.NewString()).
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
