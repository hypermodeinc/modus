/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

func GetProductVersion() string {
	return "Hypermode Runtime v" + GetVersionNumber()
}

func GetVersionNumber() string {
	return "0.x.x" // TODO
}
