/*
 * Copyright 2024 Hypermode, Inc.
 */

package plugins

import (
	"context"
	"fmt"

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

	return metadata, nil
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
