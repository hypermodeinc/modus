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
	fileIndex map[string]*Plugin
	mutex     sync.RWMutex
}

func newPluginRegistry() pluginRegistry {
	return pluginRegistry{
		// plugins:   make([]Plugin, 0),
		nameIndex: make(map[string]*Plugin),
		fileIndex: make(map[string]*Plugin),
	}
}

func (pr *pluginRegistry) Add(plugin Plugin) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	if _, ok := pr.nameIndex[plugin.Name()]; ok {
		return fmt.Errorf("plugin already exists with name %s", plugin.Name())
	}

	if _, ok := pr.fileIndex[plugin.FileName]; ok {
		return fmt.Errorf("plugin already exists with filename %s", plugin.FileName)
	}

	pr.plugins = append(pr.plugins, plugin)
	pr.nameIndex[plugin.Name()] = &plugin
	pr.fileIndex[plugin.FileName] = &plugin
	return nil
}

func (pr *pluginRegistry) Remove(plugin Plugin) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	for i, p := range pr.plugins {
		if p.Name() == plugin.Name() {
			pr.plugins[i] = pr.plugins[len(pr.plugins)-1]
			pr.plugins = pr.plugins[:len(pr.plugins)-1]
			break
		}
	}
	delete(pr.nameIndex, plugin.Name())
	delete(pr.fileIndex, plugin.FileName)
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
	if ok {
		return *plugin, true
	} else {
		return Plugin{}, false
	}
}

func (pr *pluginRegistry) GetByFile(filename string) (Plugin, bool) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	plugin, ok := pr.fileIndex[filename]
	if ok {
		return *plugin, true
	} else {
		return Plugin{}, false
	}
}
