/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package functions

import (
	"context"
	"sync"

	"github.com/hypermodeinc/modus/runtime/app"
)

type FunctionsLoadedCallback = func(ctx context.Context)

var functionsLoadedCallbacks []FunctionsLoadedCallback
var eventsMutex = sync.RWMutex{}

func RegisterFunctionsLoadedCallback(callback FunctionsLoadedCallback) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	functionsLoadedCallbacks = append(functionsLoadedCallbacks, callback)
}

func triggerFunctionsLoaded(ctx context.Context) {
	if ctx.Err() != nil || app.IsShuttingDown() {
		return
	}
	eventsMutex.RLock()
	defer eventsMutex.RUnlock()
	for _, callback := range functionsLoadedCallbacks {
		callback(ctx)
	}
}
