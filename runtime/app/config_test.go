/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package app_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
)

func TestIsDevEnvironment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expected    bool
	}{
		{"Environment is dev", "dev", true},
		{"Environment is development", "development", true},
		{"Environment is prod", "prod", false},
		{"Environment is test", "test", false},
		{"Environment is empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := app.NewAppConfig().WithEnvironment(tt.environment)
			if got := cfg.IsDevEnvironment(); got != tt.expected {
				t.Errorf("IsDevEnvironment() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewAppConfig(t *testing.T) {
	cfg := app.NewAppConfig()

	if cfg.Port() != 8686 {
		t.Errorf("Expected port to be 8686, but got %d", cfg.Port())
	}

	if cfg.RefreshInterval() != time.Second*5 {
		t.Errorf("Expected refresh interval to be 5s, but got %s", cfg.RefreshInterval())
	}

	if cfg.Environment() != "prod" {
		t.Errorf("Expected environment to be prod, but got %s", cfg.Environment())
	}

	if cfg.IsDevEnvironment() {
		t.Errorf("Expected IsDevEnvironment to be false")
	}
}

func TestCreateAppConfig(t *testing.T) {
	tests := []struct {
		name     string
		vars     map[string]string
		args     []string
		expected *app.AppConfig
	}{
		{
			name:     "default values",
			expected: app.NewAppConfig(),
		},
		{
			name: "custom values",
			vars: map[string]string{
				"MODUS_ENV": "dev",
			},
			args: []string{
				"-appPath=/path/to/app",
				"-port=9090",
				"-useAwsStorage=true",
				"-s3bucket=my-bucket",
				"-s3path=my-path",
				"-refresh=10s",
				"-jsonlogs=true",
			},
			expected: app.NewAppConfig().
				WithEnvironment("dev").
				WithPort(9090).
				WithAppPath("/path/to/app").
				WithS3Storage("my-bucket", "my-path").
				WithRefreshInterval(10 * time.Second).
				WithJsonLogging(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closer := applySettings(tt.vars, tt.args)
			t.Cleanup(closer)

			cfg := app.CreateAppConfig()

			if !reflect.DeepEqual(cfg, tt.expected) {
				t.Errorf("Expected config to be %v, but got %v", tt.expected, cfg)
			}
		})
	}
}

func applySettings(vars map[string]string, args []string) func() {
	originalVars := map[string]string{}
	for name, value := range vars {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalVars[name] = originalValue
		}
		os.Setenv(name, value)
	}

	originalArgs := os.Args
	os.Args = append([]string{os.Args[0]}, args...)

	return func() {
		for name := range vars {
			if value, ok := originalVars[name]; ok {
				os.Setenv(name, value)
			} else {
				os.Unsetenv(name)
			}
		}
		os.Args = originalArgs
	}
}
