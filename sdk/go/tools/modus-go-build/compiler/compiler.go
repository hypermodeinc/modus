/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/config"

	"github.com/hashicorp/go-version"
)

const minTinyGoVersion = "0.33.0"

func Compile(config *config.Config) error {

	tinygoVersion, err := getCompilerVersion(config)
	if err != nil {
		return err
	}

	args := []string{"build"}
	args = append(args, "-target", "wasip1")
	args = append(args, "-o", filepath.Join(config.OutputDir, config.WasmFileName))

	// WASI "reactor mode" (-buildmode=c-shared) is required for TinyGo 0.35.0 and later.
	// Otherwise, the _start function runs and immediately exits before any function can execute.
	// This also switches the startup function to _initialize instead of _start, so the Modus runtime
	// needs to match.
	if tinygoVersion.GreaterThanOrEqual(version.Must(version.NewVersion("0.35.0"))) {
		args = append(args, "-buildmode", "c-shared")
	}

	// disable the asyncify scheduler until we better understand how to use it
	args = append(args, "-scheduler", "none")

	args = append(args, config.CompilerOptions...)
	args = append(args, ".")

	cmd := exec.Command(config.CompilerPath, args...)
	cmd.Dir = config.SourceDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	return cmd.Run()
}

func Validate(config *config.Config) error {
	if !strings.HasPrefix(filepath.Base(config.CompilerPath), "tinygo") {
		return fmt.Errorf("compiler must be TinyGo")
	}

	ver, err := getCompilerVersion(config)
	if err != nil {
		return err
	}

	minVer, _ := version.NewVersion(minTinyGoVersion)
	if ver.LessThan(minVer) {
		return fmt.Errorf("found TinyGo version %s, but version %s or later is required", ver, minVer)
	}

	return nil
}

func getCompilerVersion(config *config.Config) (*version.Version, error) {
	cmd := exec.Command(config.CompilerPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(output), " ")
	if len(parts) < 3 {
		compiler := filepath.Base(config.CompilerPath)
		return nil, fmt.Errorf("unexpected output from '%s version': %s", compiler, output)
	}

	return version.NewVersion(parts[2])
}
