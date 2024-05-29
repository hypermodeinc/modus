/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import (
	"fmt"
	"os"
	"os/user"
)

const devEnvironmentName = "dev"

func GetEnvironmentName() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return devEnvironmentName
	}
	return env
}

func IsDevEnvironment() bool {
	return GetEnvironmentName() == devEnvironmentName
}

func GetNamespace() (string, error) {

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
