/*
 * Copyright 2024 Hypermode, Inc.
 */

package metadata

import (
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/hypermodeinc/modus/runtime/plugins/metadata/legacy/v1"
)

const timeFormat = "2006-01-02T15:04:05.000Z"

func getPluginMetadata_v1(wasmCustomSections map[string][]byte) (*Metadata, error) {
	metadataJson, found := wasmCustomSections["hypermode_meta"]
	if !found {
		return nil, ErrMetadataNotFound
	}

	md := v1.Metadata{}
	err := json.Unmarshal(metadataJson, &md)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plugin metadata: %w", err)
	}

	return metadataV1toV2(&md), nil
}

func metadataV1toV2(m *v1.Metadata) *Metadata {

	// legacy support for the deprecated "library" field
	// (functions-as before v0.10.0)
	sdk := m.SDK
	if sdk == "" {
		sdk = strings.TrimPrefix(m.Library, "@hypermode/")
	}

	// convert the v1 metadata to v2
	res := Metadata{
		Plugin:    m.Plugin,
		SDK:       sdk,
		BuildId:   m.BuildId,
		BuildTime: m.BuildTime.UTC().Format(timeFormat),
		GitRepo:   m.GitRepo,
		GitCommit: m.GitCommit,
		FnExports: make(FunctionMap, len(m.Functions)),
		FnImports: make(FunctionMap),
		Types:     make(TypeMap, len(m.Types)),
	}

	for _, f := range m.Functions {
		fn := &Function{Name: f.Name}

		fn.Parameters = make([]*Parameter, len(f.Parameters))
		for i, p := range f.Parameters {
			fn.Parameters[i] = &Parameter{
				Name:     p.Name,
				Type:     p.Type.Path,
				Default:  p.Default,
				Optional: p.Optional, // deprecated
			}
		}

		if f.ReturnType.Name != "" && f.ReturnType.Name != "void" {
			fn.Results = []*Result{{Type: f.ReturnType.Path}}
		}

		res.FnExports[f.Name] = fn
	}

	for _, t := range m.Types {
		td := &TypeDefinition{
			Name:   t.Path,
			Id:     t.Id,
			Fields: make([]*Field, len(t.Fields)),
		}

		for i, f := range t.Fields {
			td.Fields[i] = &Field{
				Name: f.Name,
				Type: f.Type.Path,
			}
		}

		res.Types[td.Name] = td
	}

	return &res
}
