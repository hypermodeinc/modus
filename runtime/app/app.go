/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package app

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var mu = &sync.RWMutex{}
var config *AppConfig = NewAppConfig()
var shuttingDown = false

func Initialize() {
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

// IsDevEnvironment returns true if the application is running in a development environment.
func IsDevEnvironment() bool {
	return Config().IsDevEnvironment()
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

func ModusHomeDir() string {
	modusHome := os.Getenv("MODUS_HOME")
	if modusHome == "" {
		userHome, _ := os.UserHomeDir()
		modusHome = filepath.Join(userHome, ".modus")
	}
	return modusHome
}

func KubernetesNamespace() (string, bool) {
	return kubernetesNamespace()
}

var kubernetesNamespace = sync.OnceValues(func() (string, bool) {
	if ns := os.Getenv("NAMESPACE"); ns != "" {
		return ns, true
	}
	if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		return strings.TrimSpace(string(data)), true
	}
	return "", false
})
