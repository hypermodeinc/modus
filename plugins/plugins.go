/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"hmruntime/plugins/metadata"

	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Module   wazero.CompiledModule
	Metadata metadata.Metadata
	FileName string
	Types    map[string]metadata.TypeDefinition
}

func (p *Plugin) NameAndVersion() (name string, version string) {
	return p.Metadata.NameAndVersion()
}

func (p *Plugin) Name() string {
	return p.Metadata.Name()
}

func (p *Plugin) Version() string {
	return p.Metadata.Version()
}

func (p *Plugin) BuildId() string {
	return p.Metadata.BuildId
}

func (p *Plugin) Language() metadata.PluginLanguage {
	return p.Metadata.Language()
}
