/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrMetadataNotFound = fmt.Errorf("no metadata found in plugin")

func GetMetadata(wasmCustomSections map[string][]byte) (*Metadata, error) {
	ver, err := getPluginMetadataVersion(wasmCustomSections)
	if err != nil {
		return nil, err
	}

	switch ver {
	case MetadataVersion: // current version
		return getPluginMetadata(wasmCustomSections)
	default:
		return nil, fmt.Errorf("unsupported plugin metadata version: %d", ver)
	}
}

func getPluginMetadataVersion(wasmCustomSections map[string][]byte) (byte, error) {
	verData, found := wasmCustomSections["modus_metadata_version"]
	if !found || len(verData) != 1 {
		return 0, errors.New("failed to parse plugin metadata version")
	}

	return verData[0], nil
}

func getPluginMetadata(wasmCustomSections map[string][]byte) (*Metadata, error) {
	metadataJson, found := wasmCustomSections["modus_metadata"]
	if !found {
		return nil, ErrMetadataNotFound
	}

	md := &Metadata{}
	err := json.Unmarshal(metadataJson, &md)
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
