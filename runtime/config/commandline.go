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
	"flag"
	"fmt"
	"os"
	"time"
)

var Port int
var ModelHost string
var AppPath string
var UseAwsStorage bool
var S3Bucket string
var S3Path string
var RefreshInterval time.Duration
var UseJsonLogging bool

func parseCommandLineFlags() {
	flag.StringVar(&AppPath, "appPath", "", "REQUIRED - The path to the Modus app to load and run.")
	flag.IntVar(&Port, "port", 8686, "The HTTP port to listen on.")
	flag.StringVar(&ModelHost, "modelHost", "", "The base DNS of the host endpoint to the model server.")
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
