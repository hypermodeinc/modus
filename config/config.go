/*
 * Copyright 2024 Hypermode, Inc.
 */

package config

func Initialize() {
	parseCommandLineFlags()
	setEnvironmentName()
	setNamespace()
}
