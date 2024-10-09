/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package codegen

import (
	"bytes"
	"os"
	"path/filepath"
)

const pre_file = "hyp_pre_generated.go"
const post_file = "hyp_post_generated.go"

var allFiles = []string{pre_file, post_file}

func cleanup(dir string) error {
	for _, f := range allFiles {
		err := os.Remove(filepath.Join(dir, f))
		if err != nil && !os.IsNotExist(err) {
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
