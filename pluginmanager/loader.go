/*
 * Copyright 2024 Hypermode, Inc.
 */

package pluginmanager

import (
	"context"
	"fmt"

	"hmruntime/db"
	"hmruntime/functions"
	"hmruntime/languages"
	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/plugins/metadata"
	"hmruntime/storage"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	"github.com/tetratelabs/wazero"
)

func monitorPlugins(ctx context.Context) {
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
	md, err := metadata.GetMetadata(ctx, cm)
	if err == metadata.ErrMetadataNotFound {
		logger.Error(ctx).
			Bool("user_visible", true).
			Msg("Metadata not found.  Please recompile your plugin using the latest version of the Hypermode Functions library.")
		return err
	} else if err != nil {
		return err
	}

	// Make the plugin object.
	plugin := makePlugin(ctx, cm, filename, md)

	// Write the plugin info to the database.
	// Note, this may update the ID if a plugin with the same BuildID is in the db already.
	db.WritePluginInfo(ctx, plugin)

	// Register the plugin.
	registry.AddOrUpdate(plugin)

	// Log the details of the loaded plugin.
	logPluginLoaded(ctx, plugin)

	// Trigger the plugin loaded event.
	err = triggerPluginLoaded(ctx, md)

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

func makePlugin(ctx context.Context, cm wazero.CompiledModule, filename string, md *metadata.Metadata) *plugins.Plugin {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	return &plugins.Plugin{
		Id:       utils.GenerateUUIDv7(),
		Module:   cm,
		Metadata: md,
		FileName: filename,
		Language: languages.GetLanguageForSDK(md.SDK),
	}
}

func logPluginLoaded(ctx context.Context, plugin *plugins.Plugin) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	evt := logger.Info(ctx)
	evt.Str("filename", plugin.FileName)

	md := plugin.Metadata

	if md.Plugin != "" {
		name, version := md.NameAndVersion()
		evt.Str("plugin", name)
		if version != "" {
			// The version is optional.
			evt.Str("plugin_version", version)
		}
	}

	if plugin.Language != nil {
		evt.Str("language", plugin.Language.Name())
	}

	if md.BuildId != "" {
		evt.Str("build_id", md.BuildId)
	}

	if md.BuildTime != "" {
		evt.Str("build_ts", md.BuildTime)
	}

	if md.SDK != "" {
		name, version := md.SdkNameAndVersion()
		evt.Str("sdk", name)
		evt.Str("sdk_version", version)
	}

	if md.GitRepo != "" {
		evt.Str("git_repo", md.GitRepo)
	}

	if md.GitCommit != "" {
		evt.Str("git_commit", md.GitCommit)
	}

	evt.Msg("Loaded plugin.")
}

func unloadPlugin(ctx context.Context, filename string) error {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	p := registry.GetByFile(filename)
	if p == nil {
		return fmt.Errorf("plugin not found: %s", filename)
	}

	logger.Info(ctx).
		Str("plugin", p.Name()).
		Str("build_id", p.BuildId()).
		Msg("Unloading plugin.")

	registry.Remove(p)
	return p.Module.Close(ctx)
}
