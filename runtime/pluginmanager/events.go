/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package pluginmanager

import (
	"context"
	"sync"

	"github.com/hypermodeinc/modus/runtime/plugins/metadata"
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
