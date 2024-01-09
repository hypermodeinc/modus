/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

func loadPlugins(ctx context.Context) error {
	entries, err := os.ReadDir(*pluginsPath)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {

		// Determine if the entry represents a plugin.
		var pluginName string
		entryName := entry.Name()
		if entry.IsDir() {
			pluginName = entryName
			path := fmt.Sprintf("%s/%s/build/debug.wasm", *pluginsPath, pluginName)
			if _, err := os.Stat(path); err != nil {
				continue
			}
		} else if strings.HasSuffix(entryName, ".wasm") {
			pluginName = strings.TrimSuffix(entryName, ".wasm")
		} else {
			continue
		}

		// Load the plugin
		err := loadPluginModule(ctx, pluginName)
		if err != nil {
			log.Printf("Failed to load plugin '%s': %v\n", pluginName, err)
		}
	}

	return nil
}

func watchPluginDirectory(ctx context.Context) error {
	w := watcher.New()
	w.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile(`^.+\.wasm$`), false))

	go func() {
		for {
			select {
			case evt := <-w.Event:

				pluginName, err := getPluginNameFromPath(evt.Path)
				if err != nil {
					log.Printf("failed to get plugin name: %v\n", err)
				}
				if pluginName == "" {
					continue
				}

				switch evt.Op {
				case watcher.Create, watcher.Write:
					err = loadPluginModule(ctx, pluginName)
					if err != nil {
						log.Printf("failed to load plugin: %v\n", err)
					}
				case watcher.Remove:
					err = unloadPluginModule(ctx, pluginName)
					if err != nil {
						log.Printf("failed to unload plugin: %v\n", err)
					}
				}

				// Signal that we need to register functions
				register <- true

			case err := <-w.Error:
				log.Fatalf("failure while watching plugin directory: %v\n", err)
			case <-w.Closed:
				return
			case <-ctx.Done():
				w.Close()
				return
			}
		}
	}()

	err := w.AddRecursive(*pluginsPath)
	if err != nil {
		return fmt.Errorf("failed to watch plugins directory: %w", err)
	}

	go func() {
		err = w.Start(time.Second * 1)
		if err != nil {
			log.Fatalf("failed to start file watcher: %v\n", err)
		}
	}()

	return nil
}

func getPathForPlugin(name string) (string, error) {

	// Normally the plugin will be directly in the plugins directory, by filename.
	path := *pluginsPath + "/" + name + ".wasm"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// For local development, the plugin will be in a subdirectory and we'll use the debug.wasm file.
	path = *pluginsPath + "/" + name + "/build/debug.wasm"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("compiled wasm file not found for plugin '%s'", name)
}

func getPluginNameFromPath(path string) (string, error) {
	if !strings.HasSuffix(path, ".wasm") {
		return "", fmt.Errorf("path does not point to a wasm file: %s", path)
	}

	parts := strings.Split(path, "/")

	// For local development
	if strings.HasSuffix(path, "/build/debug.wasm") {
		return parts[len(parts)-3], nil
	} else if strings.HasSuffix(path, "/build/release.wasm") {
		return "", nil
	}

	return strings.TrimSuffix(parts[len(parts)-1], ".wasm"), nil
}
