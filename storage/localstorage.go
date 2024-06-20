/*
 * Copyright 2024 Hypermode, Inc.
 */

package storage

import (
	"context"
	"fmt"
	"hmruntime/config"
	"hmruntime/logger"
	"os"
	"path/filepath"
	"strings"
)

type localStorageProvider struct {
}

func (stg *localStorageProvider) initialize(ctx context.Context) {
	if config.StoragePath == "" {
		logger.Fatal(ctx).Msg("A storage path is required when using local storage.  Exiting.")
	}

	if _, err := os.Stat(config.StoragePath); os.IsNotExist(err) {
		logger.Info(ctx).
			Str("path", config.StoragePath).
			Msg("Creating local storage directory.")
		err := os.MkdirAll(config.StoragePath, 0755)
		if err != nil {
			logger.Fatal(ctx).Err(err).
				Msg("Failed to create local storage directory.  Exiting.")
		}
	} else {
		logger.Info(ctx).
			Str("path", config.StoragePath).
			Msg("Found local storage directory.")
	}
}

func (stg *localStorageProvider) listFiles(ctx context.Context, extension string) ([]FileInfo, error) {
	entries, err := os.ReadDir(config.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in storage directory: %w", err)
	}

	var files = make([]FileInfo, 0, len(entries))
	for _, entry := range entries {

		if entry.IsDir() || !strings.HasSuffix(entry.Name(), extension) {
			continue
		}

		info, err := entry.Info()
		if err == nil {
			files = append(files, FileInfo{
				Name:         entry.Name(),
				LastModified: info.ModTime(),
			})
		}
	}

	return files, nil
}

func (stg *localStorageProvider) getFileContents(ctx context.Context, name string) ([]byte, error) {
	path := filepath.Join(config.StoragePath, name)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of file %s from local storage: %w", name, err)
	}

	return content, nil
}

func (stg *localStorageProvider) writeFile(ctx context.Context, name string, contents []byte) error {
	path := filepath.Join(config.StoragePath, name)
	err := os.WriteFile(path, contents, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s to local storage: %w", name, err)
	}

	return nil
}
