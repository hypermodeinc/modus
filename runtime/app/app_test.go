/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package app_test

import (
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/hypermodeinc/modus/runtime/app"
)

func TestIsShuttingDown(t *testing.T) {
	if app.IsShuttingDown() {
		t.Errorf("Expected initial state to be not shutting down")
	}

	app.SetShuttingDown()

	if !app.IsShuttingDown() {
		t.Errorf("Expected state to be shutting down")
	}
}

func TestSetConfig(t *testing.T) {
	initialConfig := app.Config()
	if initialConfig == nil {
		t.Errorf("Expected initial config to be non-nil")
	}

	newConfig := &app.AppConfig{}
	app.SetConfig(newConfig)

	if app.Config() != newConfig {
		t.Errorf("Expected config to be updated")
	}
}

func TestForceColor(t *testing.T) {
	if !color.NoColor {
		t.Errorf("Expected NoColor to be true")
	}

	os.Setenv("FORCE_COLOR", "1")
	app.SetOutputColorMode()

	if color.NoColor {
		t.Errorf("Expected NoColor to be false")
	}
}
