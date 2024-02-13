/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

import "flag"

var Port int
var DgraphUrl string
var PluginsPath string
var NoReload bool
var S3Bucket string

func ParseCommandLineFlags() {
	flag.IntVar(&Port, "port", 8686, "The HTTP port to listen on.")
	flag.StringVar(&DgraphUrl, "dgraph", "http://localhost:8080", "The Dgraph url to connect to.")
	flag.StringVar(&PluginsPath, "plugins", "./plugins", "The path to the plugins directory.")
	flag.StringVar(&PluginsPath, "plugin", "./plugins", "alias for -plugins")
	flag.BoolVar(&NoReload, "noreload", false, "Disable automatic plugin reloading.")
	flag.StringVar(&S3Bucket, "s3bucket", "", "The S3 bucket to use for plugin storage.")

	flag.Parse()
}
