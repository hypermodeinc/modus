/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
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
