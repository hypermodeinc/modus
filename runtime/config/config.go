/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"os"

	"github.com/fatih/color"
)

func Initialize() {
	forceColor := os.Getenv("FORCE_COLOR")
	if forceColor != "" && forceColor != "0" {
		color.NoColor = false
	}

	parseCommandLineFlags()
	readEnvironmentVariables()
}
