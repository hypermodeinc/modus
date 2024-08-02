/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"hmruntime/languages"
	"hmruntime/plugins/metadata"

	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Id       string
	Module   wazero.CompiledModule
	Metadata *metadata.Metadata
	FileName string
	Language languages.Language
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
