/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"

	"github.com/hypermodeinc/modusgraph"
)

var GlobalModusDbEngine *modusgraph.Engine

func InitModusDb(ctx context.Context) {
	if !useModusDB() {
		return
	}

	var dataDir string
	appPath := app.Config().AppPath()
	if filepath.Base(appPath) == "build" {
		// this keeps the data directory outside of the build directory
		dataDir = filepath.Join(appPath, "..", ".modusdb")
		addToGitIgnore(ctx, filepath.Dir(appPath), ".modusdb/")
	} else {
		dataDir = filepath.Join(appPath, ".modusdb")
	}

	if eng, err := modusgraph.NewEngine(modusgraph.NewDefaultConfig(dataDir)); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to initialize the local modusGraph database.")
	} else {
		GlobalModusDbEngine = eng
	}
}

func CloseModusDb(ctx context.Context) {
	if GlobalModusDbEngine != nil {
		GlobalModusDbEngine.Close()
	}
}

func addToGitIgnore(ctx context.Context, rootPath, contents string) {
	gitIgnorePath := filepath.Join(rootPath, ".gitignore")

	// if .gitignore file does not exist, create it and add contents to it
	if _, err := os.Stat(gitIgnorePath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(gitIgnorePath, []byte(contents+"\n"), 0644); err != nil {
			logger.Err(ctx, err).Msg("Failed to create .gitignore file.")
		}
		return
	}

	// check if contents are already in the .gitignore file
	file, err := os.Open(gitIgnorePath)
	if err != nil {
		logger.Err(ctx, err).Msg("Failed to open .gitignore file.")
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == contents {
			return // found
		}
	}

	// contents are not in the file, so append them
	file, err = os.OpenFile(gitIgnorePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Err(ctx, err).Msg("Failed to open .gitignore file.")
		return
	}
	defer file.Close()
	if _, err := file.WriteString("\n" + contents + "\n"); err != nil {
		logger.Err(ctx, err).Msg("Failed to append " + contents + " to .gitignore file.")
	}
}
