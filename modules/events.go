/*
 * Copyright 2024 Hypermode, Inc.
 */

package modules

import (
	"context"
	"hmruntime/plugins"
)

type PluginLoadedCallback = func(ctx context.Context, metadata plugins.PluginMetadata) error

var pluginLoadedCallbacks []PluginLoadedCallback

func RegisterPluginLoadedCallback(callback PluginLoadedCallback) {
	pluginLoadedCallbacks = append(pluginLoadedCallbacks, callback)
}

func triggerPluginLoaded(ctx context.Context, metadata plugins.PluginMetadata) error {
	for _, callback := range pluginLoadedCallbacks {
		err := callback(ctx, metadata)
		if err != nil {
			return err
		}
	}
	return nil
}
