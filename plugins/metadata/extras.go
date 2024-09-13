/*
 * Copyright 2024 Hypermode, Inc.
 */

package metadata

import "github.com/tetratelabs/wazero"

func GetMetadataFromCompiledModule(cm wazero.CompiledModule) (*Metadata, error) {
	customSections := getCustomSections(cm)
	return GetMetadata(customSections)
}

func getCustomSections(cm wazero.CompiledModule) map[string][]byte {
	sections := cm.CustomSections()
	data := make(map[string][]byte, len(sections))

	for _, s := range sections {
		data[s.Name()] = s.Data()
	}

	return data
}
