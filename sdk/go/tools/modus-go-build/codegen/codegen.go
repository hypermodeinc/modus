/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package codegen

import (
	"bytes"
	"os"
	"path/filepath"
)

const pre_file = "modus_pre_generated.go"
const post_file = "modus_post_generated.go"

func cleanup(dir string) error {
	globs := []string{"modus*_generated.go", "hyp*_generated.go"}
	files := []string{}

	for _, g := range globs {
		f, err := filepath.Glob(filepath.Join(dir, g))
		if err != nil {
			return err
		}
		files = append(files, f...)
	}

	for _, f := range files {
		err := os.Remove(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeBuffersToFile(filePath string, buffers ...*bytes.Buffer) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, b := range buffers {
		if _, err := b.WriteTo(f); err != nil {
			return err
		}
	}

	return nil
}
