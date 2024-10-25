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
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
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

var envVarsUpdated = false
var originalProcessEnvironmentVariables = os.Environ()

func GetEnvironmentName() string {
	return environment
}

func GetNamespace() string {
	return namespace
}

func readEnvironmentVariables() {
	environment = os.Getenv("MODUS_ENV")

	// default to prod
	if environment == "" {
		environment = "prod"
	}

	// If running in Kubernetes, also capture the namespace environment variable.
	namespace = os.Getenv("NAMESPACE")
}

func IsDevEnvironment() bool {
	// support either name (but prefer "dev")
	return environment == "dev" || environment == "development"
}

func getSupportedEnvironmentNames() []string {
	switch strings.ToLower(environment) {
	case "dev", "development":
		return []string{"dev", "development"}
	case "stage", "staging":
		return []string{"stage", "staging"}
	case "test", "testing":
		return []string{"test", "testing"}
	case "prod", "production":
		return []string{"prod", "production"}
	default:
		return []string{environment}
	}
}

func LoadEnvFiles(log *zerolog.Logger) {

	// Restore the original environment variables if necessary
	if envVarsUpdated {
		os.Clearenv()
		for _, envVar := range originalProcessEnvironmentVariables {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}

	// Load environment variables from .env file(s)
	envNames := getSupportedEnvironmentNames()
	files := make([]string, 0, len(envNames)*2+2)
	for _, envName := range envNames {
		files = append(files, ".env."+envName+".local")
		files = append(files, ".env."+envName)
	}
	files = append(files, ".env.local")
	files = append(files, ".env")

	for _, file := range files {
		path := filepath.Join(AppPath, file)
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				log.Warn().Err(err).Msgf("Failed to load %s file.", file)
			}
			envVarsUpdated = true
		}
	}
}
