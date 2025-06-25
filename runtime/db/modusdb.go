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
	"strings"
	"sync"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"
	mg "github.com/hypermodeinc/modusgraph"
)

var GlobalModusDbEngine *mg.Engine

var (
	client     mg.Client
	clientLock sync.Mutex
)

// GetClient returns the global modusGraph client
func GetClient() (mg.Client, error) {
	clientLock.Lock()
	defer clientLock.Unlock()
	if client == nil {
		return nil, errors.New("modusGraph client not initialized")
	}
	return client, nil
}

// InitModusDb initializes the modusGraph engine and creates
// a global modusgraph client.
//
// The `uri` parameter can be either a "file://" or "dgraph://"
// prefixed URI. The former indicates the engine will create or
// open db files in the path following `file://`. The latter
// indicates a Dgraph connection string to a remote Dgraph cluster.
//
// If `uri` is empty, the default modusGraph database file will be
// created in the app path.
func InitModusDb(ctx context.Context, uri string) {
	if !useModusDB() {
		return
	}

	var dataDir string
	if uri == "" {
		dataDir = getDefaultModusFilePath(ctx)
		dataDir = "file://" + dataDir
	} else {
		dataDir = uri
	}

	// if it's a file-based URI, create the directory if it doesn't exist
	if strings.HasPrefix(dataDir, "file://") {
		dir := strings.TrimPrefix(dataDir, "file://")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Fatal(ctx).Err(err).Str("directory", dir).Msg("Failed to create the directory for the local modusGraph database.")
		}
	}

	logVerbosity := 0
	if utils.DebugModeEnabled() {
		logVerbosity = 1
	}
	if utils.TraceModeEnabled() {
		logVerbosity = 3
	}

	var err error
	clientLock.Lock()
	defer clientLock.Unlock()
	client, err = mg.NewClient(dataDir,
		mg.WithAutoSchema(true),
		mg.WithLogger(NewZerologr(*logger.Get(ctx)).V(logVerbosity)))
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to initialize the local modusGraph database.")
	}
}

// CloseModusDb closes the modusGraph engine
// and clears the global client handle
func CloseModusDb(ctx context.Context) {
	clientLock.Lock()
	defer clientLock.Unlock()
	if client != nil {
		client.Close()
		client = nil
	}
}

func getDefaultModusFilePath(ctx context.Context) string {
	appPath := app.Config().AppPath()
	if filepath.Base(appPath) == "build" {
		addToGitIgnore(ctx, filepath.Dir(appPath), ".modusdb/")
		return filepath.Join(appPath, "..", ".modusdb")
	}
	return filepath.Join(appPath, ".modusdb")
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
