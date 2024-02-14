/*
 * Copyright 2023 Hypermode, Inc.
 */

package plugins

import (
	"context"
	"fmt"
	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/host"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

// Map of plugin names and etags as last retrieved from S3.
var awsPlugins map[string]string

func LoadPlugins(ctx context.Context) error {
	_, err := loadPlugins(ctx)
	return err
}

func ReloadPlugins(ctx context.Context) error {

	// Reload existing plugins
	loaded, err := loadPlugins(ctx)
	if err != nil {
		return err
	}

	// Unload any plugins that are no longer present
	for name := range host.CompiledModules {
		if !loaded[name] {
			err := unloadPluginModule(ctx, name)
			if err != nil {
				return fmt.Errorf("failed to unload plugin '%s': %w", name, err)
			}
		}
	}

	return nil
}

func loadPlugins(ctx context.Context) (map[string]bool, error) {
	var loaded = make(map[string]bool)

	if aws.UseAwsForPluginStorage() {
		plugins, err := aws.ListPlugins(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list plugins from S3: %w", err)
		}

		for plugin := range plugins {
			err := loadPluginModule(ctx, plugin)
			if err != nil {
				log.Printf("Failed to load plugin '%s': %v\n", plugin, err)
			} else {
				loaded[plugin] = true
			}
		}

		// Store the list of plugins and their etags for later comparison.
		awsPlugins = plugins

		return loaded, nil
	}

	// If the plugins path is a single plugin's base directory, load the single plugin.
	if _, err := os.Stat(config.PluginsPath + "/build/debug.wasm"); err == nil {
		pluginName := path.Base(config.PluginsPath)
		err := loadPluginModule(ctx, pluginName)
		if err != nil {
			log.Printf("Failed to load plugin '%s': %v\n", pluginName, err)
		} else {
			loaded[pluginName] = true
		}
	}

	// Otherwise, load all plugins in the plugins directory.
	entries, err := os.ReadDir(config.PluginsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {

		// Determine if the entry represents a plugin.
		var pluginName string
		entryName := entry.Name()
		if entry.IsDir() {
			pluginName = entryName
			path := fmt.Sprintf("%s/%s/build/debug.wasm", config.PluginsPath, pluginName)
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
		} else {
			loaded[pluginName] = true
		}
	}

	return loaded, nil
}

func WatchForPluginChanges(ctx context.Context) error {

	if config.NoReload {
		fmt.Println("NOTE: Automatic plugin reloading is disabled.  Restart the server to load new or modified plugins.")
		return nil
	}

	if aws.UseAwsForPluginStorage() {
		return watchStorageForPluginChanges(ctx)
	} else {
		return watchDirectoryForPluginChanges(ctx)
	}
}

func watchDirectoryForPluginChanges(ctx context.Context) error {

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
				host.RegistrationRequest <- true

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
	_, err := os.Lstat(config.PluginsPath)
	if err == nil {
		// They are, so we can watch recursively (local dev workflow).
		err = w.AddRecursive(config.PluginsPath)
	} else {
		// They are not.  Just watch the single directory (production workflow).
		err = w.Add(config.PluginsPath)
	}

	if err != nil {
		return fmt.Errorf("failed to watch plugins directory: %w", err)
	}

	go func() {
		err = w.Start(config.RefreshInterval)
		if err != nil {
			log.Fatalf("failed to start file watcher: %v\n", err)
		}
	}()

	return nil
}

func getPathForPlugin(name string) (string, error) {

	// Normally the plugin will be directly in the plugins directory, by filename.
	p := path.Join(config.PluginsPath, name+".wasm")
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}

	// For local development, the plugin will be in a subdirectory and we'll use the debug.wasm file.
	p = path.Join(config.PluginsPath, name, "build", "debug.wasm")
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}

	// Or, the plugins path might pointing to a single plugin's base directory.
	p = path.Join(config.PluginsPath, "build", "debug.wasm")
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}

	return "", fmt.Errorf("compiled wasm file not found for plugin '%s'", name)
}

func getPluginNameFromPath(p string) (string, error) {
	if !strings.HasSuffix(p, ".wasm") {
		return "", fmt.Errorf("path does not point to a wasm file: %s", p)
	}

	parts := strings.Split(p, "/")

	// For local development
	if strings.HasSuffix(p, "/build/debug.wasm") {
		return parts[len(parts)-3], nil
	} else if strings.HasSuffix(p, "/build/release.wasm") {
		return "", nil
	}

	return strings.TrimSuffix(parts[len(parts)-1], ".wasm"), nil
}

func watchStorageForPluginChanges(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(config.RefreshInterval)
		defer ticker.Stop()

		for {
			plugins, err := aws.ListPlugins(ctx)
			if err != nil {
				// Don't stop watching. We'll just try again on the next cycle.
				log.Println(err)
				continue
			}

			var changed = false

			// Load/reload any new or modified plugins
			for name, etag := range plugins {
				if awsPlugins[name] != etag {
					err := loadPluginModule(ctx, name)
					if err != nil {
						log.Println(err)
					}
					awsPlugins[name] = etag
					changed = true
				}
			}

			// Unload any plugins that are no longer present
			for name := range awsPlugins {
				if _, found := plugins[name]; !found {
					err := unloadPluginModule(ctx, name)
					if err != nil {
						log.Println(err)
					}
					delete(awsPlugins, name)
					changed = true
				}
			}

			// If anything changed, signal that we need to register functions
			if changed {
				host.RegistrationRequest <- true
			}

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
