/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package storage

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"

	"github.com/gofrs/flock"
)

type localStorageProvider struct {
	appPath string
}

func (stg *localStorageProvider) initialize(ctx context.Context) {

	stg.appPath = app.Config().AppPath()
	if stg.appPath == "" {
		logger.Fatal(ctx).Msg("The -appPath command line argument is required.  Exiting.")
	}

	if _, err := os.Stat(stg.appPath); os.IsNotExist(err) {
		logger.Info(ctx).
			Str("path", stg.appPath).
			Msg("Creating app directory.")
		err := os.MkdirAll(stg.appPath, 0755)
		if err != nil {
			logger.Fatal(ctx).Err(err).
				Msg("Failed to create local app directory.  Exiting.")
		}
	} else {
		logger.Info(ctx).
			Str("path", stg.appPath).
			Msg("Using local app directory.")
	}
}

func (stg *localStorageProvider) listFiles(ctx context.Context, patterns ...string) ([]FileInfo, error) {
	entries, err := os.ReadDir(stg.appPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in storage directory: %w", err)
	}

	var files = make([]FileInfo, 0, len(entries))
	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		matched := false
		for _, pattern := range patterns {
			if match, err := path.Match(pattern, filename); err == nil && match {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		info, err := entry.Info()
		if err == nil {
			files = append(files, FileInfo{
				Name:         filename,
				LastModified: info.ModTime(),
			})
		}
	}

	return files, nil
}

func (stg *localStorageProvider) getFileContents(ctx context.Context, name string) (content []byte, err error) {
	path := filepath.Join(stg.appPath, name)

	// Acquire a read lock on the file to prevent reading a file that is still being written to.
	// For example, this can easily happen when using `modus dev` and the user is editing the manifest file.

	lock := flock.New(path)
	if _, e := lock.TryRLockContext(ctx, 100*time.Millisecond); e != nil {
		return nil, fmt.Errorf("failed to acquire read lock on file %s: %w", name, e)
	}
	defer func() {
		if e := lock.Unlock(); e != nil && err == nil {
			err = fmt.Errorf("failed to release read lock on file %s: %w", name, e)
		}
	}()

	content, err = os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of file %s from local storage: %w", name, err)
	}

	return content, nil
}
