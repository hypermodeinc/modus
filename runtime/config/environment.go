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

const devEnvironmentName = "dev"

var environment string
var namespace string

func GetEnvironmentName() string {
	return environment
}

func setEnvironmentName() {
	environment = os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = devEnvironmentName
	}
}

func IsDevEnvironment() bool {
	return environment == devEnvironmentName
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
