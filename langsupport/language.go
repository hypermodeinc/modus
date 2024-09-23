/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

import (
	"hypruntime/plugins/metadata"

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
