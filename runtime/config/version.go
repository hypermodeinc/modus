/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package config

func GetProductVersion() string {
	return "Hypermode Runtime " + GetVersionNumber()
}

func GetVersionNumber() string {
	// If you see an undefined error on the following line, you need to run "go generate" in the runtime directory.
	// You can also use "make generate", "make build", or "make run", or launch in VS Code.

	return version
}
