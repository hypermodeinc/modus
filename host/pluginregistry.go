/*
 * Copyright 2024 Hypermode, Inc.
 */

package host

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
)

// Global, thread-safe registry of all plugins loaded by the host
var Plugins = newPluginRegistry()

type pluginRegistry struct {
	plugins   []Plugin
	nameIndex map[string]*Plugin
	pathIndex map[string]*Plugin
	mutex     sync.RWMutex
}

func newPluginRegistry() pluginRegistry {
	return pluginRegistry{
		// plugins:   make([]Plugin, 0),
		nameIndex: make(map[string]*Plugin),
		pathIndex: make(map[string]*Plugin),
	}
}

func (pr *pluginRegistry) Add(plugin Plugin) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	if _, ok := pr.nameIndex[plugin.Name()]; ok {
		return fmt.Errorf("plugin already exists with name %s", plugin.Name())
	}

	if _, ok := pr.pathIndex[plugin.FilePath]; ok {
		return fmt.Errorf("plugin already exists with path %s", plugin.FilePath)
	}

	pr.plugins = append(pr.plugins, plugin)
	pr.nameIndex[plugin.Name()] = &plugin
	pr.pathIndex[plugin.FilePath] = &plugin
	return nil
}

func (pr *pluginRegistry) Remove(plugin Plugin) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	for i, p := range pr.plugins {
		if p == plugin {
			pr.plugins[i] = pr.plugins[len(pr.plugins)-1]
			pr.plugins = pr.plugins[:len(pr.plugins)-1]
			break
		}
	}
	delete(pr.nameIndex, plugin.Name())
	delete(pr.pathIndex, plugin.FilePath)
}

func (pr *pluginRegistry) GetAll() []Plugin {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	plugins := pr.plugins
	slices.SortFunc(plugins, func(a, b Plugin) int {
		return cmp.Compare(a.Name(), b.Name())
	})
	return plugins
}

func (pr *pluginRegistry) GetByName(name string) (Plugin, bool) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	plugin, ok := pr.nameIndex[name]
	return *plugin, ok
}

func (pr *pluginRegistry) GetByPath(path string) (Plugin, bool) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	plugin, ok := pr.pathIndex[path]
	return *plugin, ok
}
