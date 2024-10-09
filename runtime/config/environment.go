/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"fmt"
	"os"
	"os/user"
)

/*

DESIGN NOTES:

- The MODUS_ENV environment variable is used to determine the environment name.
- We prefer to use short names, "prod", "stage", "dev", etc, but the actual value is arbitrary, so you can use longer names if you prefer.
- If it is not set, the default environment name is "prod".  This is a safe-by-default approach.
- It is preferable to actually set the MODUS_ENV to the appropriate environment when running the application.
- During development, the Modus CLI will set the MODUS_ENV to "dev" automatically.
- The "dev" environment is special in several ways, such as relaxed security requirements, and omitting certain telemetry.
- There is nothing special about "prod", other than it is the default.
- You can also use "stage", "test", etc, as needed - but they will behave like "prod".  The only difference is the name returned by the health endpoint, logs, and telemetry.

*/

var environment string
var namespace string

func GetEnvironmentName() string {
	return environment
}

func setEnvironmentName() {
	environment = os.Getenv("MODUS_ENV")

	// default to prod
	if environment == "" {
		environment = "prod"
	}
}

func IsDevEnvironment() bool {
	// support either name (but prefer "dev")
	return environment == "dev" || environment == "development"
}

func GetNamespace() string {
	return namespace
}

func setNamespace() {
	var err error
	namespace, err = getNamespaceFromOS()
	if err != nil {
		// We don't have our logger yet, so just log to stderr.
		fmt.Fprintf(os.Stderr, "Error getting namespace: %v\n", err)
		os.Exit(1)
	}
}

func getNamespaceFromOS() (string, error) {

	// In development, we'll use "dev/<username>" in lieu of the namespace.
	if IsDevEnvironment() {
		user, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("could not get current user from the os: %w", err)
		}
		return "dev/" + user.Username, nil
	}

	// Otherwise, we'll use the NAMESPACE environment variable, which is required.
	ns := os.Getenv("NAMESPACE")
	if ns == "" {
		return "", fmt.Errorf("NAMESPACE environment variable is not set")
	}

	return ns, nil
}
