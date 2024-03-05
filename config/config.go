/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import (
	"flag"
	"time"
)

var Port int
var DgraphUrl string
var PluginsPath string
var NoReload bool
var S3Bucket string
var RefreshInterval time.Duration
var UseJsonLogging bool

func ParseCommandLineFlags() {
	flag.IntVar(&Port, "port", 8686, "The HTTP port to listen on.")
	flag.StringVar(&DgraphUrl, "dgraph", "http://localhost:8080", "The Dgraph url to connect to.")
	flag.StringVar(&PluginsPath, "plugins", "", "The path to the plugins directory.")
	flag.StringVar(&PluginsPath, "plugin", "", "alias for -plugins")
	flag.BoolVar(&NoReload, "noreload", false, "Disable automatic plugin reloading.")
	flag.StringVar(&S3Bucket, "s3bucket", "", "The S3 bucket to use, if using AWS for plugin storage.")
	flag.DurationVar(&RefreshInterval, "refresh", time.Second*5, "The refresh interval to check for plugins and schema changes.")
	flag.BoolVar(&UseJsonLogging, "jsonlogs", false, "Use JSON format for logging.")

	flag.Parse()
}
