/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package envfiles

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"

	"github.com/joho/godotenv"
)

var mu sync.Mutex
var envVarsUpdated = false
var originalProcessEnvironmentVariables = os.Environ()

// Allow the env files to use short or long names.
func getSupportedEnvironmentNames() []string {
	environment := app.Config().Environment()
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

func LoadEnvFiles(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	// Restore the original environment variables if necessary
	if envVarsUpdated {
		os.Clearenv()
		for _, envVar := range originalProcessEnvironmentVariables {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
		envVarsUpdated = false
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
		path := filepath.Join(app.Config().AppPath(), file)
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				logger.Warn(ctx).Err(err).Msgf("Failed to load %s file.", file)
			}
			envVarsUpdated = true
		}
	}

	return triggerEnvFilesLoaded(ctx)
}
