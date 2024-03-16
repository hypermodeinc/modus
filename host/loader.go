/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/logger"

	"github.com/radovskyb/watcher"
)

// Map of plugin names and etags as last retrieved from S3.
var awsPluginFiles map[string]string

// Map of json files and etags as last retrieved from S3.
var awsJsons map[string]string

func LoadJsons(ctx context.Context) error {
	_, err := loadJsons(ctx)
	return err
}

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
	for _, plugin := range Plugins.GetAll() {
		if !loaded[plugin.Name()] {
			unloadPlugin(ctx, plugin)
		}
	}

	return nil
}

func loadJsons(ctx context.Context) (map[string]bool, error) {
	var loaded = make(map[string]bool)

	if aws.UseAwsForPluginStorage() {
		jsons, err := aws.ListJsons(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list jsons from S3: %w", err)
		}

		for json := range jsons {
			err := loadJson(ctx, json)
			if err != nil {
				logger.Err(ctx, err).
					Str("json", json).
					Msg("Failed to load json.")
			} else {
				loaded[json] = true
			}
		}

		// Store the list of jsons and their etags for later comparison.
		awsJsons = jsons

		return loaded, nil
	}

	// Otherwise, load all jsons in the plugins directory.
	entries, err := os.ReadDir(config.PluginsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {

		// Determine if the entry represents a json.
		var jsonName string
		entryName := entry.Name()
		if strings.HasSuffix(entryName, ".json") {
			jsonName = strings.TrimSuffix(entryName, ".json")
		} else {
			continue
		}

		// Load the json
		err := loadJson(ctx, jsonName)
		if err != nil {
			logger.Err(ctx, err).
				Str("json", jsonName).
				Msg("Failed to load json.")
		} else {
			loaded[jsonName] = true
		}
	}

	return loaded, nil
}

func loadPlugins(ctx context.Context) (map[string]bool, error) {
	var loaded = make(map[string]bool)

	if aws.UseAwsForPluginStorage() {
		files, err := aws.ListPluginsFiles(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list plugin files from S3: %w", err)
		}

		for filePath := range files {
			plugin, err := loadPlugin(ctx, filePath)
			if err != nil {
				logger.Err(ctx, err).
					Str("path", filePath).
					Msg("Failed to load plugin.")
			} else {
				pluginName := plugin.Name()
				loaded[pluginName] = true
			}
		}

		// Store the list of plugin files and their etags for later comparison.
		awsPluginFiles = files

		return loaded, nil
	}

	// If the plugins pluginPath is a single plugin's base directory, load from its build/debug.wasm file.
	pluginPath := path.Join(config.PluginsPath, "build", "debug.wasm")
	if _, err := os.Stat(pluginPath); err == nil {
		plugin, err := loadPlugin(ctx, pluginPath)
		if err != nil {
			logger.Err(ctx, err).Str("path", pluginPath).Msg("Failed to load plugin.")
		} else {
			pluginName := plugin.Name()
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
		var pluginPath string
		if entry.IsDir() {
			pluginPath = path.Join(config.PluginsPath, entry.Name(), "build", "debug.wasm")
			if _, err := os.Stat(pluginPath); err != nil {
				continue
			}
		} else if strings.HasSuffix(entry.Name(), ".wasm") {
			pluginPath = path.Join(config.PluginsPath, entry.Name())
		} else {
			continue
		}

		// Load the plugin
		plugin, err := loadPlugin(ctx, pluginPath)
		if err != nil {
			logger.Err(ctx, err).
				Str("path", pluginPath).
				Msg("Failed to load plugin.")
		} else {
			pluginName := plugin.Name()
			loaded[pluginName] = true
		}
	}

	return loaded, nil
}

func WatchForJsonChanges(ctx context.Context) error {
	if aws.UseAwsForPluginStorage() {
		return watchStorageForJsonChanges(ctx)
	} else {
		return watchDirectoryForHypermodeJsonChanges(ctx)
	}
}

func WatchForPluginChanges(ctx context.Context) error {

	if config.NoReload {
		logger.Warn(ctx).Msg("Automatic plugin reloading is disabled. Restart the server to load new or modified host.")
		return nil
	}

	if aws.UseAwsForPluginStorage() {
		return watchStorageForPluginChanges(ctx)
	} else {
		return watchDirectoryForPluginChanges(ctx)
	}
}

func watchDirectoryForHypermodeJsonChanges(ctx context.Context) error {
	w := watcher.New()
	w.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile(`^.+\.json$`), false))

	go func() {
		for {
			select {
			case evt := <-w.Event:

				jsonName, err := getJsonNameFromPath(evt.Path)
				if err != nil {
					logger.Err(ctx, err).Msg("Failed to get json name.")
				}
				if jsonName == "" {
					continue
				}

				switch evt.Op {
				case watcher.Create, watcher.Write:
					err := loadJson(ctx, jsonName)
					if err != nil {
						logger.Err(ctx, err).
							Str("json", jsonName).
							Msg("Failed to load json.")
					}
				case watcher.Remove:
					config.HypermodeData = config.HypermodeAppData{}
					logger.Info(ctx).Msg("hypermode.json removed.")
				}
			case err := <-w.Error:
				logger.Err(ctx, err).Msg("Failure while watching directory for hypermode.json")
			case <-w.Closed:
				return
			case <-ctx.Done():
				w.Close()
				return
			}
		}
	}()

	return nil
}

func watchDirectoryForPluginChanges(ctx context.Context) error {

	w := watcher.New()
	w.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile(`^.+\.wasm$`), false))

	go func() {
		for {
			select {
			case evt := <-w.Event:
				switch evt.Op {
				case watcher.Create, watcher.Write:
					_, err := loadPlugin(ctx, evt.Path)
					if err != nil {
						logger.Err(ctx, err).
							Str("path", evt.Path).
							Msg("Failed to load plugin.")
					}
				case watcher.Remove:
					plugin, found := Plugins.GetByPath(evt.Path)
					if found {
						unloadPlugin(ctx, plugin)
					}
				}

				// Signal that we need to register functions
				RegistrationRequest <- true

			case err := <-w.Error:
				logger.Err(ctx, err).Msg("Failure while watching plugin directory.")
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
			logger.Fatal(ctx).Err(err).Msg("Failed to start the file watcher.  Exiting.")
		}
	}()

	return nil
}

func getPathForJson(name string) (string, error) {
	p := path.Join(config.PluginsPath, name+".json")
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}

	return "", fmt.Errorf("json file not found for plugin '%s'", name)
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

func getJsonNameFromPath(p string) (string, error) {
	if !strings.HasSuffix(p, ".json") {
		return "", fmt.Errorf("path does not point to a json file: %s", p)
	}

	return strings.TrimSuffix(path.Base(p), ".json"), nil
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

func watchStorageForJsonChanges(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(config.RefreshInterval)
		defer ticker.Stop()

		for {
			jsons, err := aws.ListJsons(ctx)
			if err != nil {
				// Don't stop watching. We'll just try again on the next cycle.
				logger.Err(ctx, err).Msg("Failed to list jsons from S3.")
				continue
			}
			// Load/reload any new or modified jsons
			for name, etag := range jsons {
				if awsJsons[name] != etag {
					err := loadJson(ctx, name)
					if err != nil {
						logger.Err(ctx, err).
							Str("json", name).
							Msg("Failed to load hypermode.json.")
					}
					awsJsons[name] = etag
				}
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

func watchStorageForPluginChanges(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(config.RefreshInterval)
		defer ticker.Stop()

		for {
			pluginFiles, err := aws.ListPluginsFiles(ctx)
			if err != nil {
				// Don't stop watching. We'll just try again on the next cycle.
				logger.Err(ctx, err).Msg("Failed to list plugins from S3.")
				continue
			}

			var changed = false

			// Load/reload any new or modified plugins
			for filename, etag := range pluginFiles {
				if awsPluginFiles[filename] != etag {
					_, err := loadPlugin(ctx, filename)
					if err != nil {
						logger.Err(ctx, err).
							Str("path", filename).
							Msg("Failed to load plugin.")
					}
					awsPluginFiles[filename] = etag
					changed = true
				}
			}

			// Unload any plugins that are no longer present
			for filename := range awsPluginFiles {
				if _, found := pluginFiles[filename]; !found {
					plugin, found := Plugins.GetByPath(filename)
					if found {
						unloadPlugin(ctx, plugin)
					}
					delete(awsPluginFiles, filename)
					changed = true
				}
			}

			// If anything changed, signal that we need to register functions
			if changed {
				RegistrationRequest <- true
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
