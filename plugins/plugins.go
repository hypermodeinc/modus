/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"context"
	"hmruntime/languages"
	"hmruntime/plugins/metadata"
	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Id       string
	Module   wazero.CompiledModule
	Metadata *metadata.Metadata
	FileName string
	Language languages.Language
}

func NewPlugin(cm wazero.CompiledModule, filename string, md *metadata.Metadata) *Plugin {
	return &Plugin{
		Id:       utils.GenerateUUIDv7(),
		Module:   cm,
		Metadata: md,
		FileName: filename,
		Language: languages.GetLanguageForSDK(md.SDK),
	}
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

func GetPlugin(ctx context.Context) *Plugin {
	if p, ok := ctx.Value(utils.PluginContextKey).(*Plugin); ok {
		return p
	}

	return nil
}
