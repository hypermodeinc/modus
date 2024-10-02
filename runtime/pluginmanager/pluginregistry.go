/*
 * Copyright 2024 Hypermode, Inc.
 */

package pluginmanager

import (
	"cmp"
	"slices"
	"sync"

	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/utils"
)

// thread-safe globalPluginRegistry of all plugins that are loaded
var globalPluginRegistry = &pluginRegistry{
	idRevIndex: make(map[*plugins.Plugin]string),
	idIndex:    make(map[string]*plugins.Plugin),
	nameIndex:  make(map[string]*plugins.Plugin),
	fileIndex:  make(map[string]*plugins.Plugin),
}

func GetRegisteredPlugins() []*plugins.Plugin {
	return globalPluginRegistry.GetAll()
}

type pluginRegistry struct {
	idRevIndex map[*plugins.Plugin]string
	idIndex    map[string]*plugins.Plugin
	nameIndex  map[string]*plugins.Plugin
	fileIndex  map[string]*plugins.Plugin
	mutex      sync.RWMutex
}

func (pr *pluginRegistry) AddOrUpdate(plugin *plugins.Plugin) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	// only one plugin per name is allowed
	name := plugin.Name()
	if existing, found := pr.nameIndex[name]; found {
		delete(pr.idRevIndex, existing)
		delete(pr.idIndex, existing.Id)
		delete(pr.nameIndex, name)
		delete(pr.fileIndex, existing.FileName)
	}

	pr.idRevIndex[plugin] = plugin.Id
	pr.idIndex[plugin.Id] = plugin
	pr.nameIndex[name] = plugin
	pr.fileIndex[plugin.FileName] = plugin
}

func (pr *pluginRegistry) Remove(plugin *plugins.Plugin) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	id, found := pr.idRevIndex[plugin]
	if !found {
		return
	}

	delete(pr.idRevIndex, plugin)
	delete(pr.idIndex, id)
	delete(pr.nameIndex, plugin.Name())
	delete(pr.fileIndex, plugin.FileName)
}

func (pr *pluginRegistry) GetAll() []*plugins.Plugin {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	result := utils.MapKeys(pr.idRevIndex)
	slices.SortFunc(result, func(a, b *plugins.Plugin) int {
		return cmp.Compare(a.Name(), b.Name())
	})
	return result
}

func (pr *pluginRegistry) GetId(plugin *plugins.Plugin) string {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	if id, found := pr.idRevIndex[plugin]; found {
		return id
	}

	return ""
}

func (pr *pluginRegistry) GetById(id string) *plugins.Plugin {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	if plugin, found := pr.idIndex[id]; found {
		return plugin
	}

	return nil
}

func (pr *pluginRegistry) GetByName(name string) *plugins.Plugin {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	if plugin, found := pr.nameIndex[name]; found {
		return plugin
	}

	return nil
}

func (pr *pluginRegistry) GetByFile(filename string) *plugins.Plugin {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	if plugin, found := pr.fileIndex[filename]; found {
		return plugin
	}

	return nil
}
