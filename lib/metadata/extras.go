/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package metadata

import (
	"github.com/hypermodeinc/modus/lib/wasmextractor"
)

func GetMetadataFromCompiledModule(cm []byte) (*Metadata, error) {
	customSections, err := getCustomSections(cm)
	if err != nil {
		return nil, err
	}
	return GetMetadata(customSections)
}

func getCustomSections(cm []byte) (map[string][]byte, error) {
	info, err := wasmextractor.ExtractWasmInfo(cm)
	if err != nil {
		return nil, err
	}
	return info.CustomSections, nil
}
