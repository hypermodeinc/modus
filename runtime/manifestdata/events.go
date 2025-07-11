/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifestdata

import (
	"context"
	"sync"

	"github.com/hypermodeinc/modus/runtime/app"
)

type ManifestLoadedCallback = func(ctx context.Context) error

var manifestLoadedCallbacks []ManifestLoadedCallback
var eventsMutex = sync.RWMutex{}

func RegisterManifestLoadedCallback(callback ManifestLoadedCallback) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	manifestLoadedCallbacks = append(manifestLoadedCallbacks, callback)
}

func triggerManifestLoaded(ctx context.Context) error {
	if ctx.Err() != nil || app.IsShuttingDown() {
		return nil
	}

	eventsMutex.RLock()
	defer eventsMutex.RUnlock()

	for _, callback := range manifestLoadedCallbacks {
		err := callback(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
