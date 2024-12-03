/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package app

import (
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/fatih/color"
)

var mu = &sync.RWMutex{}
var config *AppConfig
var shuttingDown = false

func init() {
	// Set the global color mode
	SetOutputColorMode()

	// Create the the app configuration
	mu.Lock()
	defer mu.Unlock()
	config = CreateAppConfig()
}

// SetOutputColorMode applies the FORCE_COLOR environment variable to override the default color mode.
func SetOutputColorMode() {
	forceColor := os.Getenv("FORCE_COLOR")
	if forceColor != "" && forceColor != "0" {
		color.NoColor = false
	}
}

// Config returns the global app configuration.
func Config() *AppConfig {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

// SetConfig sets the global app configuration.
// This is typically only called in tests.
func SetConfig(c *AppConfig) {
	mu.Lock()
	defer mu.Unlock()
	config = c
}

// IsShuttingDown returns true if the application is in the process of a graceful shutdown.
func IsShuttingDown() bool {
	mu.RLock()
	defer mu.RUnlock()
	return shuttingDown
}

// SetShuttingDown sets the application to a shutting down state during a graceful shutdown.
func SetShuttingDown() {
	mu.Lock()
	defer mu.Unlock()
	shuttingDown = true
}

// GetRootSourcePath returns the root path of the source code.
// It is used to trim the paths in stack traces when included in telemetry.
func GetRootSourcePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return path.Dir(path.Dir(filename)) + "/"
}
