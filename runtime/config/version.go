/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

func GetProductVersion() string {
	return "Hypermode Runtime " + GetVersionNumber()
}

func GetVersionNumber() string {
	return version
}
