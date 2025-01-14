/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package langsupport

import (
	"github.com/hypermodeinc/modus/lib/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

type Language interface {
	Name() string
	TypeInfo() LanguageTypeInfo
	NewPlanner(md *metadata.Metadata) Planner
	NewWasmAdapter(mod wasm.Module) WasmAdapter
}

func NewLanguage(name string, typeInfo LanguageTypeInfo, plannerFactory func(*metadata.Metadata) Planner, waFactory func(wasm.Module) WasmAdapter) Language {
	return &language{name, typeInfo, plannerFactory, waFactory}
}

type language struct {
	name           string
	typeInfo       LanguageTypeInfo
	plannerFactory func(*metadata.Metadata) Planner
	waFactory      func(wasm.Module) WasmAdapter
}

func (l *language) Name() string {
	return l.name
}

func (l *language) TypeInfo() LanguageTypeInfo {
	return l.typeInfo
}

func (l *language) NewPlanner(md *metadata.Metadata) Planner {
	return l.plannerFactory(md)
}

func (l *language) NewWasmAdapter(mod wasm.Module) WasmAdapter {
	return l.waFactory(mod)
}
