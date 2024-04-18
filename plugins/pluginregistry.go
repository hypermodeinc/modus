/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"cmp"
	"slices"
	"sync"
)

type pluginRegistry struct {
	plugins   []Plugin
	nameIndex map[string]*Plugin
	fileIndex map[string]*Plugin
	mutex     sync.RWMutex
}

func NewPluginRegistry() pluginRegistry {
	return pluginRegistry{
		nameIndex: make(map[string]*Plugin),
		fileIndex: make(map[string]*Plugin),
	}
}

func (pr *pluginRegistry) AddOrUpdate(plugin Plugin) {
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
