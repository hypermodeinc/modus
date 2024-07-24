/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import "hmruntime/plugins"

type ModelInfo struct {
	Name     string
	FullName string
}

func (m *ModelInfo) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "ModelInfo",
		Path: "~lib/@hypermode/models-as/index/ModelInfo",
	}
}
