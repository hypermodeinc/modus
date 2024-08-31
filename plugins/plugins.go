/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"context"

	"hypruntime/langsupport"
	"hypruntime/languages"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"

	"github.com/tetratelabs/wazero"
)

type Plugin struct {
	Id       string
	Module   wazero.CompiledModule
	Metadata *metadata.Metadata
	FileName string
	Language langsupport.Language
	Planner  langsupport.Planner
}

func NewPlugin(cm wazero.CompiledModule, filename string, md *metadata.Metadata) *Plugin {

	language := languages.GetLanguageForSDK(md.SDK)
	planner := language.NewPlanner(md)

	return &Plugin{
		Id:       utils.GenerateUUIDv7(),
		Module:   cm,
		Metadata: md,
		FileName: filename,
		Language: language,
		Planner:  planner,
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
