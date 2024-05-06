/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import "os"

func GetProductVersion() string {
	return "Hypermode Runtime " + GetVersionNumber()
}

func GetVersionNumber() string {
	return version
}

func GetEnvironmentName() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return "dev"
	}
	return env
}
