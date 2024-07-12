/*
 * Copyright 2024 Hypermode, Inc.
 */

package pluginmanager

import (
	"context"
	"fmt"
	"time"

	"hmruntime/functions"
	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/storage"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	"github.com/tetratelabs/wazero"
)

func MonitorPlugins(ctx context.Context) {
	loadPluginFile := func(fi storage.FileInfo) error {
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
			plugins := registry.GetAll()
			functions.RegisterFunctions(ctx, plugins)
		}
	}
	sm.Start(ctx)
}

func loadPlugin(ctx context.Context, filename string) error {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// Load the binary content of the plugin.
	bytes, err := storage.GetFileContents(ctx, filename)
	if err != nil {
		return err
	}

	// Compile the plugin into a module
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

	// Trigger the plugin loaded event.
	err = triggerPluginLoaded(ctx, metadata)

	return err
}

func compileModule(ctx context.Context, bytes []byte) (wazero.CompiledModule, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	cm, err := wasmhost.RuntimeInstance.CompileModule(ctx, bytes)
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
	registry.AddOrUpdate(plugin)

	return plugin
}

func logPluginLoaded(ctx context.Context, plugin plugins.Plugin) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	evt := logger.Info(ctx)
	evt.Str("filename", plugin.FileName)

	metadata := plugin.Metadata

	if metadata.Plugin != "" {
		name, version := metadata.NameAndVersion()
		evt.Str("plugin", name)
		if version != "" {
			// The version is optional.
			evt.Str("plugin_version", version)
		}
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

	if metadata.SDK != "" {
		name, version := metadata.SdkNameAndVersion()
		evt.Str("sdk", name)
		evt.Str("sdk_version", version)
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
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	p, ok := registry.GetByFile(filename)
	if !ok {
		return fmt.Errorf("plugin not found: %s", filename)
	}
	logger.Info(ctx).
		Str("plugin", p.Name()).
		Str("build_id", p.BuildId()).
		Msg("Unloading plugin.")

	registry.Remove(p)
	return (*p.Module).Close(ctx)
}
