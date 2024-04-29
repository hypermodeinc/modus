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

	// Load the binary content of the plugin.
	bytes, err := storage.GetFileContents(ctx, filename)
	if err != nil {
		return err
	}

	// Compile the plugin into a module.
	cm, err := RuntimeInstance.CompileModule(ctx, bytes)
	if err != nil {
		return fmt.Errorf("failed to compile the plugin: %w", err)
	}

	// Get the metadata for the plugin.
	metadata, err := plugins.GetPluginMetadata(&cm)
	if err == plugins.ErrPluginMetadataNotFound {
		logger.Error(ctx).
			Msg("Metadata not found.  Please recompile your plugin using the latest version of the Hypermode Functions library.")
		return err
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
	Plugins.AddOrUpdate(plugin)

	// Log the details of the loaded plugin.
	logPluginLoaded(ctx, plugin)

	// Notify the callback that a plugin has been loaded.
	err = pluginLoaded(ctx, metadata)
	if err != nil {
		return err
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
