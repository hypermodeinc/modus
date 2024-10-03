/*
 * Copyright 2024 Hypermode, Inc.
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hypermodeAI/functions-go/tools/hypbuild/codegen"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/compiler"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/config"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/metagen"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/modinfo"
	"github.com/hypermodeAI/functions-go/tools/hypbuild/wasm"
	"github.com/hypermodeAI/manifest"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		exitWithError("Error", err)
	}

	if err := compiler.Validate(config); err != nil {
		exitWithError("Error", err)
	}

	mod, err := modinfo.CollectModuleInfo(config)
	if err != nil {
		exitWithError("Error", err)
	}

	color.NoColor = !isatty.IsTerminal(os.Stdout.Fd())

	metagen.WriteLogo()

	if err := codegen.PreProcess(config); err != nil {
		exitWithError("Error while pre-processing source files", err)
	}

	msg := fmt.Sprintf("\nBuilding %s ...", config.WasmFileName)
	fmt.Printf("%s\n\n", msg)

	if err := compiler.Compile(config, false); err != nil {
		exitWithError("Error building wasm", err)
	}

	meta, err := metagen.GenerateMetadata(config, mod)
	if err != nil {
		exitWithError("Error generating metadata", err)
	}

	if err := codegen.PostProcess(config, meta); err != nil {
		exitWithError("Error while post-processing source files", err)
	}

	if err := compiler.Compile(config, true); err != nil {
		exitWithError("Error building wasm", err)
	}

	if err := wasm.WriteMetadata(config, meta); err != nil {
		exitWithError("Error writing metadata", err)
	}

	if err := validateAndCopyManifestToOutput(config); err != nil {
		exitWithError("Manifest error", err)
	}

	// for dramatic effect
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("\033[2A\033[%dC\U0001F389\n\n", len(msg))
		time.Sleep(250 * time.Millisecond)
	}

	metagen.LogToConsole(meta)
}

func exitWithError(msg string, err error) {
	fmt.Fprintf(os.Stderr, msg+": %v\n", err)
	os.Exit(1)
}

func validateAndCopyManifestToOutput(config *config.Config) error {
	manifestFile := filepath.Join(config.SourceDir, "hypermode.json")
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

	outFile := filepath.Join(config.OutputDir, "hypermode.json")
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
