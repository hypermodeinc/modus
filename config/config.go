/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

// command line flag variables
var DgraphUrl string
var PluginsPath string
var NoReload bool

// Channel used to signal that registration is needed
var Register chan bool = make(chan bool)

// channel and flag used to signal the HTTP server
var ServerReady chan bool = make(chan bool)
var ServerWaiting = true
