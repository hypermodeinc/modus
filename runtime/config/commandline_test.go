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
	"flag"
	"os"
	"testing"
	"time"
)

func TestParseCommandLineFlags(t *testing.T) {
	tests := []struct {
		name                    string
		args                    []string
		expectedPort            int
		expectedAppPath         string
		expectedUseAwsStorage   bool
		expectedS3Bucket        string
		expectedS3Path          string
		expectedRefreshInterval time.Duration
		expectedUseJsonLogging  bool
	}{
		{
			name:                    "default values",
			args:                    []string{},
			expectedPort:            8686,
			expectedAppPath:         "",
			expectedUseAwsStorage:   false,
			expectedS3Bucket:        "",
			expectedS3Path:          "",
			expectedRefreshInterval: time.Second * 5,
			expectedUseJsonLogging:  false,
		},
		{
			name: "custom values",
			args: []string{
				"-appPath=/path/to/app",
				"-port=9090",
				"-useAwsStorage=true",
				"-s3bucket=my-bucket",
				"-s3path=my-path",
				"-refresh=10s",
				"-jsonlogs=true",
			},
			expectedPort:            9090,
			expectedAppPath:         "/path/to/app",
			expectedUseAwsStorage:   true,
			expectedS3Bucket:        "my-bucket",
			expectedS3Path:          "my-path",
			expectedRefreshInterval: 10 * time.Second,
			expectedUseJsonLogging:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags and variables
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			Port = 0
			AppPath = ""
			UseAwsStorage = false
			S3Bucket = ""
			S3Path = ""
			RefreshInterval = 0
			UseJsonLogging = false

			// Set command line arguments
			os.Args = append([]string{os.Args[0]}, tt.args...)

			// Parse flags
			parseCommandLineFlags()

			// Check values
			if Port != tt.expectedPort {
				t.Errorf("expected Port %d, got %d", tt.expectedPort, Port)
			}
			if AppPath != tt.expectedAppPath {
				t.Errorf("expected AppPath %s, got %s", tt.expectedAppPath, AppPath)
			}
			if UseAwsStorage != tt.expectedUseAwsStorage {
				t.Errorf("expected UseAwsStorage %v, got %v", tt.expectedUseAwsStorage, UseAwsStorage)
			}
			if S3Bucket != tt.expectedS3Bucket {
				t.Errorf("expected S3Bucket %s, got %s", tt.expectedS3Bucket, S3Bucket)
			}
			if S3Path != tt.expectedS3Path {
				t.Errorf("expected S3Path %s, got %s", tt.expectedS3Path, S3Path)
			}
			if RefreshInterval != tt.expectedRefreshInterval {
				t.Errorf("expected RefreshInterval %v, got %v", tt.expectedRefreshInterval, RefreshInterval)
			}
			if UseJsonLogging != tt.expectedUseJsonLogging {
				t.Errorf("expected UseJsonLogging %v, got %v", tt.expectedUseJsonLogging, UseJsonLogging)
			}
		})
	}
}
