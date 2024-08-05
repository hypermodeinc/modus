/*
 * Copyright 2024 Hypermode, Inc.
 */

package metadata

import (
	"context"
	"fmt"
	"strings"

	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
)

var ErrPluginMetadataNotFound = fmt.Errorf("no metadata found in plugin")

func GetPluginMetadata(ctx context.Context, cm wazero.CompiledModule) (*Metadata, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	metadataJson, found := getCustomSectionData(cm, "hypermode_meta")
	if !found {
		return nil, ErrPluginMetadataNotFound
	}

	md := &Metadata{}
	err := utils.JsonDeserialize(metadataJson, &md)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plugin metadata: %w", err)
	}

	augmentMetadata(md)
	return md, nil
}

func augmentMetadata(metadata *Metadata) {

	// legacy support for the deprecated "library" field
	// (functions-as before v0.10.0)
	if metadata.SDK == "" {
		metadata.SDK = strings.TrimPrefix(metadata.Library, "@hypermode/")
	}

	// Copy the language from the metadata to the types and functions.
	// Set the nullable flag along the way.
	lang := metadata.Language()
	for i, t := range metadata.Types {
		for j, field := range t.Fields {
			field.Type.Language = lang
			field.Type.Nullable = isNullable(field.Type, lang)
			t.Fields[j] = field
		}
		metadata.Types[i] = t
	}
	for i, fn := range metadata.Functions {
		for j, param := range fn.Parameters {
			param.Type.Language = lang
			param.Type.Nullable = isNullable(param.Type, lang)
			fn.Parameters[j] = param
		}
		fn.ReturnType.Language = lang
		fn.ReturnType.Nullable = isNullable(fn.ReturnType, lang)
		metadata.Functions[i] = fn
	}
}

func isNullable(t TypeInfo, lang PluginLanguage) bool {
	switch lang {
	case AssemblyScript:
		return strings.HasSuffix(t.Path, "|null")
	default:
		panic(fmt.Sprintf("unsupported language: %v", lang))
	}
}

func getCustomSectionData(cm wazero.CompiledModule, name string) (data []byte, found bool) {
	for _, sec := range cm.CustomSections() {
		if sec.Name() == name {
			data = sec.Data()
			found = true
			break
		}
	}
	return data, found
}
