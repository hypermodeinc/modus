/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package wasm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hypermodeAI/functions-go/tools/hypbuild/config"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/metadata"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/utils"
	"github.com/hypermodeAI/wasmextractor"
)

func WriteMetadata(config *config.Config, meta *metadata.Metadata) error {
	metaJson, err := utils.JsonSerialize(meta, false)
	if err != nil {
		return err
	}

	wasmFilePath := filepath.Join(config.OutputDir, config.WasmFileName)
	return writeCustomSections(wasmFilePath, map[string][]byte{
		"hypermode_version": {metadata.MetadataVersion},
		"hypermode_meta":    metaJson,
	})
}

func writeCustomSections(wasmFilePath string, customSections map[string][]byte) error {
	wasmBytes, err := getWasmBytes(wasmFilePath)
	if err != nil {
		return err
	}

	tmpFilePath := wasmFilePath[:len(wasmFilePath)-5] + "_tmp.wasm"
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer os.Remove(tmpFilePath)

	// Preamble
	_, _ = tmpFile.Write(wasmBytes[:8])
	offset := 8

	// Sections
	for offset < len(wasmBytes) {
		secStart := offset

		sectionID := wasmBytes[offset]
		offset++

		size, n := binary.Uvarint(wasmBytes[offset:])
		offset += n

		// Skip existing custom section with the same names as the new ones
		if sectionID == 0 {
			nameLen, n := binary.Uvarint(wasmBytes[offset:])
			nameData := wasmBytes[offset+n : offset+n+int(nameLen)]
			if _, ok := customSections[string(nameData)]; ok {
				offset += int(size)
				continue
			}
		}

		offset += int(size)
		_, _ = tmpFile.Write(wasmBytes[secStart:offset])
	}

	// Add new custom sections
	for name, data := range customSections {
		secNameBytes := []byte(name)
		secNameLen := makeUvarint(len(secNameBytes))
		payloadSize := len(secNameLen) + len(secNameBytes) + len(data)
		secSize := makeUvarint(payloadSize)

		_, _ = tmpFile.Write([]byte{0})
		_, _ = tmpFile.Write(secSize)
		_, _ = tmpFile.Write(secNameLen)
		_, _ = tmpFile.Write(secNameBytes)
		_, _ = tmpFile.Write(data)
	}

	tmpFile.Close()
	return os.Rename(tmpFilePath, wasmFilePath)
}

func makeUvarint(x int) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(x))
	return buf[:n]
}

func getWasmBytes(wasmFilePath string) ([]byte, error) {
	wasmBytes, err := os.ReadFile(wasmFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading wasm file: %v", err)
	}

	magic := []byte{0x00, 0x61, 0x73, 0x6D} // "\0asm"
	if len(wasmBytes) < 8 || !bytes.Equal(wasmBytes[:4], magic) {
		return nil, fmt.Errorf("invalid wasm file")
	}

	if binary.LittleEndian.Uint32(wasmBytes[4:8]) != 1 {
		return nil, fmt.Errorf("unsupported wasm version")
	}

	return wasmBytes, nil
}

type WasmFunctions struct {
	Exports map[string]any
	Imports map[string]any
}

func GetWasmFunctions(wasmFilePath string) (*WasmFunctions, error) {
	wasmBytes, err := getWasmBytes(wasmFilePath)
	if err != nil {
		return nil, err
	}

	info, err := wasmextractor.ExtractWasmInfo(wasmBytes)
	if err != nil {
		return nil, err
	}

	result := WasmFunctions{}

	shouldIgnore := func(name string) bool {
		ignorePrefixes := []string{"wasi", "env", "runtime", "syscall", "gojs"}
		for _, prefix := range ignorePrefixes {
			if strings.HasPrefix(name, prefix) {
				return true
			}
		}
		return false
	}

	result.Imports = make(map[string]any, len(info.Imports))
	for _, item := range info.Imports {
		if item.Kind == wasmextractor.WasmFunction && !shouldIgnore(item.Name) {
			result.Imports[item.Name] = nil
		}
	}

	result.Exports = make(map[string]any, len(info.Exports))
	for _, item := range info.Exports {
		if item.Kind == wasmextractor.WasmFunction && !shouldIgnore(item.Name) {
			result.Exports[item.Name] = nil
		}
	}

	return &result, nil
}
