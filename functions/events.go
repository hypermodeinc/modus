/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import "context"

type FunctionsLoadedCallback = func(ctx context.Context)

var functionsLoadedCallbacks []FunctionsLoadedCallback

func RegisterFunctionsLoadedCallback(callback FunctionsLoadedCallback) {
	functionsLoadedCallbacks = append(functionsLoadedCallbacks, callback)
}

func triggerFunctionsLoaded(ctx context.Context) {
	for _, callback := range functionsLoadedCallbacks {
		callback(ctx)
	}
}
