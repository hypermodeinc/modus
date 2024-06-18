/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import (
	"fmt"
	"os"
	"os/user"
)

const DevEnvironmentName = "dev"

var environment string
var namespace string

func GetEnvironmentName() string {
	return environment
}

func setEnvironmentName() {
	environment = os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = DevEnvironmentName
	}
}

func IsDevEnvironment() bool {
	return environment == DevEnvironmentName
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
