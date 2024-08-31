/*
 * Copyright 2024 Hypermode, Inc.
 */

package metadata

import (
	"context"
	"errors"
	"fmt"

	v1 "hypruntime/plugins/metadata/legacy/v1"
	"hypruntime/utils"

	"github.com/tetratelabs/wazero"
)

var ErrMetadataNotFound = fmt.Errorf("no metadata found in plugin")

func GetMetadata(ctx context.Context, cm wazero.CompiledModule) (*Metadata, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	ver, err := getPluginMetadataVersion(cm)
	if err != nil {
		return nil, err
	}

	switch ver {
	case 1:
		return getPluginMetadata_v1(ctx, cm)
	case 2:
		return getPluginMetadata_v2(ctx, cm)
	default:
		return nil, fmt.Errorf("unsupported plugin metadata version: %d", ver)
	}
}

func getPluginMetadataVersion(cm wazero.CompiledModule) (byte, error) {
	verData, found := getCustomSectionData(cm, "hypermode_version")
	if !found {
		return 1, nil // version 1 did not have a version section
	}

	if len(verData) != 1 {
		return 0, errors.New("failed to parse plugin metadata version")
	}

	return verData[0], nil
}

func getPluginMetadata_v1(ctx context.Context, cm wazero.CompiledModule) (*Metadata, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	metadataJson, found := getCustomSectionData(cm, "hypermode_meta")
	if !found {
		return nil, ErrMetadataNotFound
	}

	md := v1.Metadata{}
	err := utils.JsonDeserialize(metadataJson, &md)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plugin metadata: %w", err)
	}

	return metadataV1toV2(&md), nil
}

func getPluginMetadata_v2(ctx context.Context, cm wazero.CompiledModule) (*Metadata, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	metadataJson, found := getCustomSectionData(cm, "hypermode_meta")
	if !found {
		return nil, ErrMetadataNotFound
	}

	md := &Metadata{}
	err := utils.JsonDeserialize(metadataJson, &md)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plugin metadata: %w", err)
	}

	for name, fn := range md.FnExports {
		fn.Name = name
		md.FnExports[name] = fn
	}

	for name, fn := range md.FnImports {
		fn.Name = name
		md.FnImports[name] = fn
	}

	for name, typ := range md.Types {
		typ.Name = name
		md.Types[name] = typ
	}

	return md, nil
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
