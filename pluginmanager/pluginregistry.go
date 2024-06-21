/*
 * Copyright 2024 Hypermode, Inc.
 */

package pluginmanager

import (
	"cmp"
	"hmruntime/plugins"
	"slices"
	"sync"
)

// thread-safe registry of all plugins that are loaded
var registry = pluginRegistry{
	nameIndex: make(map[string]*plugins.Plugin),
	fileIndex: make(map[string]*plugins.Plugin),
}

type pluginRegistry struct {
	plugins   []plugins.Plugin
	nameIndex map[string]*plugins.Plugin
	fileIndex map[string]*plugins.Plugin
	mutex     sync.RWMutex
}

func (pr *pluginRegistry) AddOrUpdate(plugin plugins.Plugin) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	_, found := pr.nameIndex[plugin.Name()]
	if found {
		for i, p := range pr.plugins {
			if p.Name() == plugin.Name() {
				pr.plugins[i] = plugin
				break
			}
		}
	} else {
		pr.plugins = append(pr.plugins, plugin)
	}

	pr.nameIndex[plugin.Name()] = &plugin
	pr.fileIndex[plugin.FileName] = &plugin
}

func (pr *pluginRegistry) Remove(plugin plugins.Plugin) {
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

func (pr *pluginRegistry) GetAll() []plugins.Plugin {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	result := pr.plugins
	slices.SortFunc(result, func(a, b plugins.Plugin) int {
		return cmp.Compare(a.Name(), b.Name())
	})
	return result
}

func (pr *pluginRegistry) GetByName(name string) (plugins.Plugin, bool) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	plugin, ok := pr.nameIndex[name]
	if ok {
		return *plugin, true
	} else {
		return plugins.Plugin{}, false
	}
}

func (pr *pluginRegistry) GetByFile(filename string) (plugins.Plugin, bool) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	plugin, ok := pr.fileIndex[filename]
	if ok {
		return *plugin, true
	} else {
		return plugins.Plugin{}, false
	}
}
