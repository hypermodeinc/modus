/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import (
	"context"
	"fmt"

	"hmruntime/logger"
	"hmruntime/storage"
)

func LoadJsons(ctx context.Context) error {
	_, err := loadJsons(ctx)
	return err
}

func ReloadPlugins(ctx context.Context) error {

	// Reload existing plugins
	loaded, err := loadPlugins(ctx)
	if err != nil {
		return err
	}

	// Unload any plugins that are no longer present
	for _, plugin := range Plugins.GetAll() {
		if !loaded[plugin.Name()] {
			unloadPlugin(ctx, plugin)
		}
	}

	return nil
}

func loadJsons(ctx context.Context) (map[string]bool, error) {
	var loaded = make(map[string]bool)

	files, err := storage.ListFiles(ctx, ".json")
	if err != nil {
		return nil, fmt.Errorf("failed to list JSON files: %w", err)
	}

	for _, file := range files {
		err := loadJson(ctx, file.Name)
		if err != nil {
			logger.Err(ctx, err).
				Str("filename", file.Name).
				Msg("Failed to load JSON file.")
		} else {
			loaded[file.Name] = true
		}
	}

	return loaded, nil
}

func loadPlugins(ctx context.Context) (map[string]bool, error) {
	var loaded = make(map[string]bool)

	files, err := storage.ListFiles(ctx, ".wasm")
	if err != nil {
		return nil, fmt.Errorf("failed to list plugin files: %w", err)
	}

	for _, file := range files {
		plugin, err := loadPlugin(ctx, file.Name)
		if err != nil {
			logger.Err(ctx, err).
				Str("filename", file.Name).
				Msg("Failed to load plugin file.")
		} else {
			loaded[plugin.Name()] = true
		}
	}

	return loaded, nil
}
