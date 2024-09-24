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
	IsImport() bool
	Plugin() *plugins.Plugin
	Metadata() *metadata.Function
	ExecutionPlan() langsupport.ExecutionPlan
}

func NewFunctionInfo(fnName string, plugin *plugins.Plugin, isImport bool) (FunctionInfo, bool) {

	var fnMap metadata.FunctionMap
	if isImport {
		fnMap = plugin.Metadata.FnImports
	} else {
		fnMap = plugin.Metadata.FnExports
	}

	fnMeta := fnMap[fnName]
	plan := plugin.ExecutionPlans[fnName]
	if fnMeta == nil || plan == nil {
		return nil, false
	}

	info := &functionInfo{fnName, isImport, plugin, fnMeta, plan}
	return info, true
}

type functionInfo struct {
	fnName   string
	isImport bool
	plugin   *plugins.Plugin
	fnMeta   *metadata.Function
	plan     langsupport.ExecutionPlan
}

func (f *functionInfo) Name() string                             { return f.fnName }
func (f *functionInfo) IsImport() bool                           { return f.isImport }
func (f *functionInfo) Plugin() *plugins.Plugin                  { return f.plugin }
func (f *functionInfo) Metadata() *metadata.Function             { return f.fnMeta }
func (f *functionInfo) ExecutionPlan() langsupport.ExecutionPlan { return f.plan }
