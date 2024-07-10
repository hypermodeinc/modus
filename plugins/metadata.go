/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"context"
	"fmt"
	"strings"

	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
)

var ErrPluginMetadataNotFound = fmt.Errorf("no metadata found in plugin")

func GetPluginMetadata(ctx context.Context, cm *wazero.CompiledModule) (PluginMetadata, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	metadataJson, found := getCustomSectionData(cm, "hypermode_meta")
	if !found {
		return PluginMetadata{}, ErrPluginMetadataNotFound
	}

	metadata := PluginMetadata{}
	err := utils.JsonDeserialize(metadataJson, &metadata)
	if err != nil {
		return PluginMetadata{}, fmt.Errorf("failed to parse plugin metadata: %w", err)
	}

	augmentMetadata(&metadata)
	return metadata, nil
}

func augmentMetadata(metadata *PluginMetadata) {
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

func getCustomSectionData(cm *wazero.CompiledModule, name string) (data []byte, found bool) {
	for _, sec := range (*cm).CustomSections() {
		if sec.Name() == name {
			data = sec.Data()
			found = true
			break
		}
	}
	return data, found
}
