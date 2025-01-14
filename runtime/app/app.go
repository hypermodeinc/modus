/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package app

import (
	"path"
	"runtime"
	"sync"
	"time"
)

// ShutdownTimeout is the time to wait for the server to shutdown gracefully.
const ShutdownTimeout = 5 * time.Second

var mu = &sync.RWMutex{}
var shuttingDown = false

func IsShuttingDown() bool {
	mu.RLock()
	defer mu.RUnlock()
	return shuttingDown
}

func SetShuttingDown() {
	mu.Lock()
	defer mu.Unlock()
	shuttingDown = true
}

func GetRootSourcePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return path.Join(path.Dir(filename), "../") + "/"
}
