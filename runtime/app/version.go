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
	"os/exec"
	"strings"
)

var version string

func init() {
	adjustVersion()
}

func adjustVersion() {
	// The "version" variable is set by the makefile using -ldflags when using "make build" or goreleaser.
	// If it is not set, then we are running in development mode with "go run" or "go build" without the makefile,
	// so we will describe the version from git at run time.
	if version == "" {
		version = describeVersion()
	} else if version[0] >= '0' && version[0] <= '9' {
		version = "v" + version
	}
}

func VersionNumber() string {
	return version
}

func ProductVersion() string {
	return "Modus Runtime " + VersionNumber()
}

func describeVersion() string {
	result, err := exec.Command("git", "describe", "--tags", "--always", "--match", "runtime/*").Output()
	if err != nil {
		return "(unknown)"
	}
	return strings.TrimPrefix(strings.TrimSpace(string(result)), "runtime/")
}
