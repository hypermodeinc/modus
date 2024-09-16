/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"context"
	"sync"
)

type FunctionsLoadedCallback = func(ctx context.Context) bool

var functionsLoadedCallbacks []FunctionsLoadedCallback
var eventsMutex = sync.RWMutex{}

func RegisterFunctionsLoadedCallback(callback FunctionsLoadedCallback) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	functionsLoadedCallbacks = append(functionsLoadedCallbacks, callback)
}

func triggerFunctionsLoaded(ctx context.Context) {
	eventsMutex.RLock()
	defer eventsMutex.RUnlock()
	for _, callback := range functionsLoadedCallbacks {
		callback(ctx)
	}
}
