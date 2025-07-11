/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	SourceDir       string
	CompilerPath    string
	CompilerOptions []string
	WasmFileName    string
	OutputDir       string
}

func GetConfig() (*Config, error) {
	c := &Config{}

	defaultCompilerPath, err := getCompilerDefaultPath()
	if err != nil {
		return nil, err
	}

	flag.StringVar(&c.CompilerPath, "compiler", defaultCompilerPath, "Path to the compiler to use.")
	flag.StringVar(&c.OutputDir, "output", "build", "Output directory for the generated files. Relative paths are resolved relative to the source directory.")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "modus-go-build - The build tool for Go-based Modus apps")
		fmt.Fprintln(os.Stderr, "Usage: modus-go-build [options] <source directory> [... additional compiler options ...]")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	c.SourceDir = flag.Arg(0)

	if flag.NArg() > 1 {
		c.CompilerOptions = flag.Args()[1:]
	}

	if _, err := os.Stat(c.SourceDir); err != nil {
		return nil, fmt.Errorf("error locating source directory: %w", err)
	}

	wasmFileName, err := getWasmFileName(c.SourceDir)
	if err != nil {
		return nil, fmt.Errorf("error determining output filename: %w", err)
	}
	c.WasmFileName = wasmFileName

	if !filepath.IsAbs(c.OutputDir) {
		c.OutputDir = filepath.Join(c.SourceDir, c.OutputDir)
	}

	if err := os.MkdirAll(c.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating output directory: %w", err)
	}

	return c, nil
}

func getWasmFileName(sourceDir string) (string, error) {
	sourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return "", err
	}
	return filepath.Base(sourceDir) + ".wasm", nil
}

func getCompilerDefaultPath() (string, error) {
	for _, arg := range os.Args {
		if strings.TrimLeft(arg, "-") == "compiler" {
			return "", nil
		}
	}

	path, err := exec.LookPath("tinygo")
	if err != nil {
		return "", fmt.Errorf("tinygo not found in PATH - see https://tinygo.org/getting-started/install for installation instructions")
	}

	return path, nil
}
