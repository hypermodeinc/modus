/*
 * Copyright 2024 Hypermode, Inc.
 */

package pluginmanager

import (
	"context"
	"hmruntime/plugins"
	"sync"
)

type PluginLoadedCallback = func(ctx context.Context, metadata plugins.PluginMetadata) error

var pluginLoadedCallbacks []PluginLoadedCallback
var eventsMutex = sync.RWMutex{}

func RegisterPluginLoadedCallback(callback PluginLoadedCallback) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	pluginLoadedCallbacks = append(pluginLoadedCallbacks, callback)
}

func triggerPluginLoaded(ctx context.Context, metadata plugins.PluginMetadata) error {
	eventsMutex.RLock()
	defer eventsMutex.RUnlock()
	for _, callback := range pluginLoadedCallbacks {
		err := callback(ctx, metadata)
		if err != nil {
			return err
		}
	}
	return nil
}
