/*
 * Copyright 2024 Hypermode, Inc.
 */

package pluginmanager

import (
	"context"
	"sync"

	"hmruntime/plugins/metadata"
)

type PluginLoadedCallback = func(ctx context.Context, md *metadata.Metadata) error

var pluginLoadedCallbacks []PluginLoadedCallback
var eventsMutex = sync.RWMutex{}

func RegisterPluginLoadedCallback(callback PluginLoadedCallback) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	pluginLoadedCallbacks = append(pluginLoadedCallbacks, callback)
}

func triggerPluginLoaded(ctx context.Context, md *metadata.Metadata) error {
	eventsMutex.RLock()
	defer eventsMutex.RUnlock()
	for _, callback := range pluginLoadedCallbacks {
		err := callback(ctx, md)
		if err != nil {
			return err
		}
	}
	return nil
}
