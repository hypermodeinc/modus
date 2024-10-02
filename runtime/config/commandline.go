/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var Port int
var ModelHost string
var StoragePath string
var UseAwsSecrets bool
var UseAwsStorage bool
var S3Bucket string
var S3Path string
var RefreshInterval time.Duration
var UseJsonLogging bool

func parseCommandLineFlags() {
	flag.IntVar(&Port, "port", 8686, "The HTTP port to listen on.")
	flag.StringVar(&ModelHost, "modelHost", "", "The base DNS of the host endpoint to the model server.")
	flag.StringVar(&StoragePath, "storagePath", getDefaultStoragePath(), "The path to a directory used for local storage.")
	flag.BoolVar(&UseAwsSecrets, "useAwsSecrets", false, "Use AWS Secrets Manager for API keys and other secrets.")
	flag.BoolVar(&UseAwsStorage, "useAwsStorage", false, "Use AWS S3 for storage instead of the local filesystem.")
	flag.StringVar(&S3Bucket, "s3bucket", "", "The S3 bucket to use, if using AWS storage.")
	flag.StringVar(&S3Path, "s3path", "", "The path within the S3 bucket to use, if using AWS storage.")
	flag.DurationVar(&RefreshInterval, "refresh", time.Second*5, "The refresh interval to reload any changes.")
	flag.BoolVar(&UseJsonLogging, "jsonlogs", false, "Use JSON format for logging.")

	var showVersion bool
	const versionUsage = "Show the Runtime version number and exit."
	flag.BoolVar(&showVersion, "version", false, versionUsage)
	flag.BoolVar(&showVersion, "v", false, versionUsage+" (shorthand)")

	flag.Parse()

	if showVersion {
		fmt.Println(GetProductVersion())
		os.Exit(0)
	}
}

func getDefaultStoragePath() string {

	// On Windows, the default is %APPDATA%\Hypermode
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		return filepath.Join(appData, "Hypermode")
	}

	// On Unix and macOS, the default is $HOME/.hypermode
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homedir, ".hypermode")
}
