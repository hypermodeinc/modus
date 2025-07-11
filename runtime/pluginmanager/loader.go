/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package pluginmanager

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/storage"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
)

func monitorPlugins(ctx context.Context) {
	loadPluginFile := func(fi storage.FileInfo) error {
		err := loadPlugin(ctx, fi.Name)
		if err != nil {
			const msg = "Failed to load plugin."
			sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("filename", fi.Name))
			logger.Error(ctx, err).Str("filename", fi.Name).Msg(msg)
		}
		return err
	}

	sm := storage.NewStorageMonitor("*.wasm")
	sm.Added = loadPluginFile
	sm.Modified = loadPluginFile
	sm.Removed = func(fi storage.FileInfo) error {
		err := unloadPlugin(ctx, fi.Name)
		if err != nil {
			const msg = "Failed to unload plugin."
			sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("filename", fi.Name))
			logger.Error(ctx, err).Str("filename", fi.Name).Msg(msg)
		}
		return err
	}
	sm.Start(ctx)
}

func loadPlugin(ctx context.Context, filename string) error {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	// Load the binary content of the plugin.
	bytes, err := storage.GetFileContents(ctx, filename)
	if err != nil {
		return err
	}

	// Compile the plugin into a module
	host := wasmhost.GetWasmHost(ctx)
	cm, err := host.CompileModule(ctx, bytes)
	if err != nil {
		return err
	}

	// Get the metadata for the plugin.
	md, err := metadata.GetMetadataFromWasm(bytes)
	if err == metadata.ErrMetadataNotFound {
		const msg = "Metadata not found.  Please recompile using the latest version of the Modus SDK."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Error(ctx).Bool("user_visible", true).Msg(msg)
		return err
	} else if err != nil {
		return err
	}

	// Make the plugin object.
	plugin, err := plugins.NewPlugin(ctx, cm, filename, md)
	if err != nil {
		return err
	}

	// Write the plugin info to the database.
	// Note, this may update the ID if a plugin with the same BuildID is in the db already.
	db.WritePluginInfo(ctx, plugin)

	// Register the plugin.
	globalPluginRegistry.AddOrUpdate(plugin)

	// Log the details of the loaded plugin.
	logPluginLoaded(ctx, plugin)

	// register functions in the plugin
	host.GetFunctionRegistry().RegisterAllFunctions(ctx, plugin)

	// Trigger the plugin loaded event.
	err = triggerPluginLoaded(ctx, plugin)

	return err
}

func logPluginLoaded(ctx context.Context, plugin *plugins.Plugin) {
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
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	p, found := globalPluginRegistry.GetByFile(filename)
	if !found {
		return fmt.Errorf("plugin not found: %s", filename)
	}

	logger.Info(ctx).
		Str("plugin", p.Name()).
		Str("build_id", p.BuildId()).
		Msg("Unloading plugin.")

	globalPluginRegistry.Remove(p)
	return p.Module.Close(ctx)
}
