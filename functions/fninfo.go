/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"hypruntime/langsupport"
	"hypruntime/plugins"
	"hypruntime/plugins/metadata"
)

type FunctionInfo interface {
	Name() string
	Plugin() *plugins.Plugin
	Metadata() *metadata.Function
	ExecutionPlan() langsupport.ExecutionPlan
}

func NewFunctionInfo(plugin *plugins.Plugin, plan langsupport.ExecutionPlan) FunctionInfo {
	return &functionInfo{
		plugin:        plugin,
		executionPlan: plan,
	}
}

type functionInfo struct {
	plugin        *plugins.Plugin
	executionPlan langsupport.ExecutionPlan
}

func (f *functionInfo) Name() string {
	return f.Metadata().Name
}

func (f *functionInfo) Plugin() *plugins.Plugin {
	return f.plugin
}

func (f *functionInfo) Metadata() *metadata.Function {
	return f.executionPlan.FnMetadata()
}

func (f *functionInfo) ExecutionPlan() langsupport.ExecutionPlan {
	return f.executionPlan
}
