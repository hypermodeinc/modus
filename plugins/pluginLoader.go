/*
 * Copyright 2023 Hypermode, Inc.
 */
package plugins

import (
	"context"
	"fmt"
	"hmruntime/monitor"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

// Polling interval to check for new plugins
const pluginRefreshInterval time.Duration = time.Second * 5

var PluginsPath *string

func LoadPlugins(ctx context.Context) error {

	// If the plugins path is a single plugin's base directory, load the single plugin.
	if _, err := os.Stat(*PluginsPath + "/build/debug.wasm"); err == nil {
		pluginName := path.Base(*PluginsPath)
		err := loadPluginModule(ctx, pluginName)
		if err != nil {
			log.Printf("Failed to load plugin '%s': %v\n", pluginName, err)
		}
	}

	// Otherwise, load all plugins in the plugins directory.
	entries, err := os.ReadDir(*PluginsPath)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {

		// Determine if the entry represents a plugin.
		var pluginName string
		entryName := entry.Name()
		if entry.IsDir() {
			pluginName = entryName
			path := fmt.Sprintf("%s/%s/build/debug.wasm", *PluginsPath, pluginName)
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

func WatchPluginDirectory(ctx context.Context) error {
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
				monitor.Register <- true

			case err := <-w.Error:
				log.Printf("failure while watching plugin directory: %v\n", err)
			case <-w.Closed:
				return
			case <-ctx.Done():
				w.Close()
				return
			}
		}
	}()

	// Test if symlinks are supported
	_, err := os.Lstat(*PluginsPath)
	if err == nil {
		// They are, so we can watch recursively (local dev workflow).
		err = w.AddRecursive(*PluginsPath)
	} else {
		// They are not.  Just watch the single directory (production workflow).
		err = w.Add(*PluginsPath)
	}

	if err != nil {
		return fmt.Errorf("failed to watch plugins directory: %w", err)
	}

	go func() {
		err = w.Start(pluginRefreshInterval)
		if err != nil {
			log.Fatalf("failed to start file watcher: %v\n", err)
		}
	}()

	return nil
}

func getPathForPlugin(name string) (string, error) {

	// Normally the plugin will be directly in the plugins directory, by filename.
	path := *PluginsPath + "/" + name + ".wasm"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// For local development, the plugin will be in a subdirectory and we'll use the debug.wasm file.
	path = *PluginsPath + "/" + name + "/build/debug.wasm"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// Or, the plugins path might pointing to a single plugin's base directory.
	path = *PluginsPath + "/build/debug.wasm"
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
