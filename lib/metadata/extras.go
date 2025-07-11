/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package metadata

import (
	"github.com/hypermodeinc/modus/lib/wasmextractor"
)

func GetMetadataFromWasm(wasm []byte) (*Metadata, error) {
	customSections, err := getCustomSections(wasm)
	if err != nil {
		return nil, err
	}
	return GetMetadata(customSections)
}

func getCustomSections(wasm []byte) (map[string][]byte, error) {
	info, err := wasmextractor.ExtractWasmInfo(wasm)
	if err != nil {
		return nil, err
	}
	return info.CustomSections, nil
}
