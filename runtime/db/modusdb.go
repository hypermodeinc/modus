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
	"fmt"
	"os"
	"path/filepath"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"

	"github.com/hypermodeinc/modusgraph"
)

var globalModusDbEngine *modusgraph.Engine

func InitModusDb(ctx context.Context) {
	if !useModusDB() {
		return
	}

	var dataDir string
	appPath := app.Config().AppPath()
	if filepath.Base(appPath) == "build" {
		// this keeps the data directory outside of the build directory
		dataDir = filepath.Join(appPath, "..", ".modusdb")
		if err := addToGitIgnore(ctx, filepath.Dir(appPath), ".modusdb/"); err != nil {
			const msg = "Failed to add .modusdb to .gitignore"
			sentryutils.CaptureWarning(ctx, err, msg)
			logger.Warn(ctx, err).Msg(msg)
		}
	} else {
		dataDir = filepath.Join(appPath, ".modusdb")
	}

	if eng, err := modusgraph.NewEngine(modusgraph.NewDefaultConfig(dataDir)); err != nil {
		const msg = "Failed to initialize the local modusGraph database."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Fatal(ctx, err).Msg(msg)
	} else {
		globalModusDbEngine = eng
	}
}

func CloseModusDb(ctx context.Context) {
	if globalModusDbEngine != nil {
		globalModusDbEngine.Close()
	}
}

func addToGitIgnore(ctx context.Context, rootPath, contents string) error {
	gitIgnorePath := filepath.Join(rootPath, ".gitignore")

	// if .gitignore file does not exist, create it and add contents to it
	if _, err := os.Stat(gitIgnorePath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(gitIgnorePath, []byte(contents+"\n"), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore file: %w", err)
		}
		return nil
	}

	// check if contents are already in the .gitignore file
	file, err := os.Open(gitIgnorePath)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore file: %w", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == contents {
			file.Close()
			return nil // found
		}
	}
	file.Close()

	// contents are not in the file, so append them
	file, err = os.OpenFile(gitIgnorePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore file: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString("\n" + contents + "\n"); err != nil {
		return fmt.Errorf("failed to append to .gitignore file: %w", err)
	}
	return nil
}
