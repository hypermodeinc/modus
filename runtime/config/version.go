/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

func GetProductVersion() string {
	return "Hypermode Runtime " + GetVersionNumber()
}

func GetVersionNumber() string {
	// If you see an undefined error on the following line, you need to run "go generate" in the runtime directory.
	// You can also use "make generate", "make build", or "make run", or launch in VS Code.

	return version
}
