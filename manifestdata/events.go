/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifestdata

import (
	"context"
	"sync"
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
