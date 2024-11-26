/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/codegen"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/compiler"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/config"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metagen"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/modinfo"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/wasm"
)

func main() {

	start := time.Now()
	trace := utils.IsTraceModeEnabled()
	if trace {
		log.Println("Starting build process...")
	}

	config, err := config.GetConfig()
	if err != nil {
		exitWithError("Error", err)
	}

	if trace {
		log.Println("Configuration loaded.")
	}

	if err := compiler.Validate(config); err != nil {
		exitWithError("Error", err)
	}

	if trace {
		log.Println("Configuration validated.")
	}

	mod, err := modinfo.CollectModuleInfo(config)
	if err != nil {
		exitWithError("Error", err)
	}

	if trace {
		log.Println("Module info collected.")
	}

	if err := codegen.PreProcess(config); err != nil {
		exitWithError("Error while pre-processing source files", err)
	}

	if trace {
		log.Println("Pre-processing done.")
	}

	meta, err := metagen.GenerateMetadata(config, mod)
	if err != nil {
		exitWithError("Error generating metadata", err)
	}

	if trace {
		log.Println("Metadata generated.")
	}

	if err := codegen.PostProcess(config, meta); err != nil {
		exitWithError("Error while post-processing source files", err)
	}

	if trace {
		log.Println("Post-processing done.")
	}

	if err := compiler.Compile(config); err != nil {
		exitWithError("Error building wasm", err)
	}

	if trace {
		log.Println("Wasm compiled (second pass).")
	}

	if err := wasm.WriteMetadata(config, meta); err != nil {
		exitWithError("Error writing metadata", err)
	}

	if trace {
		log.Println("Metadata written.")
	}

	if err := validateAndCopyManifestToOutput(config); err != nil {
		exitWithError("Manifest error", err)
	}

	if trace {
		log.Println("Manifest copied.")
	}

	if trace {
		log.Printf("Build completed in %.2f seconds.\n\n", time.Since(start).Seconds())
	}

	metagen.LogToConsole(meta)
}

func exitWithError(msg string, err error) {
	fmt.Fprintf(os.Stderr, msg+": %v\n", err)
	os.Exit(1)
}

func validateAndCopyManifestToOutput(config *config.Config) error {
	manifestFile := filepath.Join(config.SourceDir, "modus.json")
	if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	data, err := os.ReadFile(manifestFile)
	if err != nil {
		return err
	}

	if err := validateManifest(data); err != nil {
		return err
	}

	outFile := filepath.Join(config.OutputDir, "modus.json")
	if err := os.WriteFile(outFile, data, 0644); err != nil {
		return err
	}

	return nil
}

func validateManifest(data []byte) error {
	// Make a copy of the data to avoid modifying the original
	// TODO: this should be fixed in the manifest library
	manData := make([]byte, len(data))
	copy(manData, data)
	return manifest.ValidateManifest(manData)
}
