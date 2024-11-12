/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package envfiles

import (
	"context"
	"sync"

	"github.com/hypermodeinc/modus/runtime/app"
)

type EnvFilesLoadedCallback = func(ctx context.Context)

var manifestLoadedCallbacks []EnvFilesLoadedCallback
var eventsMutex = sync.RWMutex{}

func RegisterEnvFilesLoadedCallback(callback EnvFilesLoadedCallback) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	manifestLoadedCallbacks = append(manifestLoadedCallbacks, callback)
}

func triggerEnvFilesLoaded(ctx context.Context) error {
	if app.IsShuttingDown() {
		return nil
	}

	eventsMutex.RLock()
	defer eventsMutex.RUnlock()

	for _, callback := range manifestLoadedCallbacks {
		callback(ctx)
	}
	return nil
}
