/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package functions

import (
	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/plugins/metadata"
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
