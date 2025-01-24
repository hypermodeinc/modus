/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package app

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	environment     string
	port            int
	appPath         string
	useAwsStorage   bool
	s3Bucket        string
	s3Path          string
	refreshInterval time.Duration
	useJsonLogging  bool
}

func (c *AppConfig) Environment() string {
	return c.environment
}

func (c *AppConfig) Port() int {
	return c.port
}

func (c *AppConfig) AppPath() string {
	return c.appPath
}

func (c *AppConfig) UseAwsStorage() bool {
	return c.useAwsStorage
}

func (c *AppConfig) S3Bucket() string {
	return c.s3Bucket
}

func (c *AppConfig) S3Path() string {
	return c.s3Path
}

func (c *AppConfig) RefreshInterval() time.Duration {
	return c.refreshInterval
}

func (c *AppConfig) UseJsonLogging() bool {
	return c.useJsonLogging
}

func (c *AppConfig) IsDevEnvironment() bool {
	// support either name (but prefer "dev")
	return c.environment == "dev" || c.environment == "development"
}

func (c *AppConfig) WithEnvironment(environment string) *AppConfig {
	cfg := *c
	cfg.environment = environment
	return &cfg
}

func (c *AppConfig) WithPort(port int) *AppConfig {
	cfg := *c
	cfg.port = port
	return &cfg
}

func (c *AppConfig) WithAppPath(appPath string) *AppConfig {
	cfg := *c
	cfg.appPath = appPath
	return &cfg
}

func (c *AppConfig) WithS3Storage(s3Bucket, s3Path string) *AppConfig {
	cfg := *c
	cfg.useAwsStorage = true
	cfg.s3Bucket = s3Bucket
	cfg.s3Path = s3Path
	return &cfg
}

func (c *AppConfig) WithRefreshInterval(interval time.Duration) *AppConfig {
	cfg := *c
	cfg.refreshInterval = interval
	return &cfg
}

func (c *AppConfig) WithJsonLogging() *AppConfig {
	cfg := *c
	cfg.useJsonLogging = true
	return &cfg
}

// Creates a new AppConfig instance with default values.
func NewAppConfig() *AppConfig {
	return &AppConfig{
		port:            8686,
		environment:     "prod",
		refreshInterval: time.Second * 5,
	}
}

// Creates the app configuration from the command line flags and environment variables.
func CreateAppConfig() *AppConfig {

	cfg := NewAppConfig()

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&cfg.appPath, "appPath", cfg.appPath, "REQUIRED - The path to the Modus app to load and run.")
	fs.IntVar(&cfg.port, "port", cfg.port, "The HTTP port to listen on.")

	fs.BoolVar(&cfg.useAwsStorage, "useAwsStorage", cfg.useAwsStorage, "Use AWS S3 for storage instead of the local filesystem.")
	fs.StringVar(&cfg.s3Bucket, "s3bucket", cfg.s3Bucket, "The S3 bucket to use, if using AWS storage.")
	fs.StringVar(&cfg.s3Path, "s3path", cfg.s3Path, "The path within the S3 bucket to use, if using AWS storage.")

	fs.DurationVar(&cfg.refreshInterval, "refresh", cfg.refreshInterval, "The refresh interval to reload any changes.")
	fs.BoolVar(&cfg.useJsonLogging, "jsonlogs", cfg.useJsonLogging, "Use JSON format for logging.")

	var showVersion bool
	const versionUsage = "Show the Runtime version number and exit."
	fs.BoolVar(&showVersion, "version", false, versionUsage)
	fs.BoolVar(&showVersion, "v", false, versionUsage+" (shorthand)")

	args := make([]string, 0, len(os.Args))
	for i := 1; i < len(os.Args); i++ {
		if !strings.HasPrefix(os.Args[i], "-test.") {
			args = append(args, os.Args[i])
		}
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing command line flags: %v\n", err)
		os.Exit(1)
	}

	if showVersion {
		fmt.Println(ProductVersion())
		os.Exit(0)
	}

	if env := os.Getenv("MODUS_ENV"); env != "" {
		cfg.environment = env
	}

	return cfg
}
