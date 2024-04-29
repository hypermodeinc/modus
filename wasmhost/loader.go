/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"fmt"
	"time"

	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/storage"
	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
)

var pluginLoaded func(ctx context.Context, metadata plugins.PluginMetadata) error

// Registers a callback function that is called when a plugin is loaded.
func RegisterPluginLoadedCallback(callback func(ctx context.Context, metadata plugins.PluginMetadata) error) {
	pluginLoaded = callback
}

func MonitorPlugins() {
	loadPluginFile := func(fi storage.FileInfo) error {
		ctx := context.Background()
		err := loadPlugin(ctx, fi.Name)
		if err != nil {
			logger.Err(ctx, err).
				Str("filename", fi.Name).
				Msg("Failed to load plugin.")
		}
		return err
	}

	sm := storage.NewStorageMonitor(".wasm")
	sm.Added = loadPluginFile
	sm.Modified = loadPluginFile
	sm.Removed = func(fi storage.FileInfo) error {
		ctx := context.Background()
		err := unloadPlugin(ctx, fi.Name)
		if err != nil {
			logger.Err(ctx, err).
				Str("filename", fi.Name).
				Msg("Failed to unload plugin.")
		}
		return err
	}
	sm.Changed = func(errors []error) {
		if len(errors) == 0 {
			// Signal that we need to register functions
			RegistrationRequest <- true
		}
	}
	sm.Start()
}

func loadPlugin(ctx context.Context, filename string) error {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// Load the binary content of the plugin.
	bytes, err := storage.GetFileContents(ctx, filename)
	if err != nil {
		return err
	}

	// Compile the plugin into a module.
	cm, err := compileModule(ctx, bytes)
	if err != nil {
		return err
	}

	// Get the metadata for the plugin.
	metadata, err := plugins.GetPluginMetadata(ctx, &cm)
	if err == plugins.ErrPluginMetadataNotFound {
		logger.Error(ctx).
			Bool("user_visible", true).
			Msg("Metadata not found.  Please recompile your plugin using the latest version of the Hypermode Functions library.")
		return err
	} else if err != nil {
		return err
	}

	// Make the plugin object.
	plugin := makePlugin(ctx, &cm, filename, metadata)

	// Log the details of the loaded plugin.
	logPluginLoaded(ctx, plugin)

	// Notify the callback that a plugin has been loaded.
	err = pluginLoaded(ctx, metadata)
	if err != nil {
		return err
	}

	return nil
}

func compileModule(ctx context.Context, bytes []byte) (wazero.CompiledModule, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	cm, err := RuntimeInstance.CompileModule(ctx, bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile the plugin: %w", err)
	}

	return cm, nil
}

func makePlugin(ctx context.Context, cm *wazero.CompiledModule, filename string, metadata plugins.PluginMetadata) plugins.Plugin {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	// Store the types in a map for easy access.
	types := make(map[string]plugins.TypeDefinition, len(metadata.Types))
	for _, t := range metadata.Types {
		types[t.Path] = t
	}

	// Create and store the plugin.
	plugin := plugins.Plugin{Module: cm, Metadata: metadata, FileName: filename, Types: types}
	Plugins.AddOrUpdate(plugin)

	return plugin
}

func logPluginLoaded(ctx context.Context, plugin plugins.Plugin) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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

func unloadPlugin(ctx context.Context, filename string) error {
	p, ok := Plugins.GetByFile(filename)
	if !ok {
		return fmt.Errorf("plugin not found: %s", filename)
	}
	logger.Info(ctx).
		Str("plugin", p.Name()).
		Str("build_id", p.BuildId()).
		Msg("Unloading plugin.")

	Plugins.Remove(p)
	return (*p.Module).Close(ctx)
}
