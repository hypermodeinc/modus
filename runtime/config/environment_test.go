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
	"testing"
)

func TestEnvironmentNames(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		expectedResult string
		isDev          bool
	}{
		{
			name:           "Environment variable not set",
			envValue:       "",
			expectedResult: "prod",
			isDev:          false,
		},
		{
			name:           "Environment variable set to dev",
			envValue:       "dev",
			expectedResult: "dev",
			isDev:          true,
		},
		{
			name:           "Environment variable set to development",
			envValue:       "development",
			expectedResult: "development",
			isDev:          true,
		},
		{
			name:           "Environment variable set to stage",
			envValue:       "stage",
			expectedResult: "stage",
			isDev:          false,
		},
		{
			name:           "Environment variable set to test",
			envValue:       "test",
			expectedResult: "test",
			isDev:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("MODUS_ENV", tt.envValue)
			readEnvironmentVariables()
			result := GetEnvironmentName()
			if result != tt.expectedResult {
				t.Errorf("Expected environment to be %s, but got %s", tt.expectedResult, result)
			}
			if IsDevEnvironment() != tt.isDev {
				t.Errorf("Expected IsDevEnvironment to be %v, but got %v", tt.isDev, IsDevEnvironment())
			}
		})
	}
}
