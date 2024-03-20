/*
 * Copyright 2024 Hypermode, Inc.
 */

package storage

import (
	"context"
	"fmt"
	"hmruntime/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type localStorage struct {
}

func (stg *localStorage) initialize() {
	if config.StoragePath == "" {
		log.Fatal().Msg("A storage path is required when using local storage.  Exiting.")
	}

	if _, err := os.Stat(config.StoragePath); os.IsNotExist(err) {
		log.Info().
			Str("path", config.StoragePath).
			Msg("Creating local storage directory.")
		err := os.MkdirAll(config.StoragePath, 0755)
		if err != nil {
			log.Fatal().Err(err).
				Msg("Failed to create local storage directory.  Exiting.")
		}
	} else {
		log.Info().
			Str("path", config.StoragePath).
			Msg("Found local storage directory.")
	}
}

func (stg *localStorage) listFiles(ctx context.Context, extension string) ([]FileInfo, error) {
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

func (stg *localStorage) getFileContents(ctx context.Context, name string) ([]byte, error) {
	path := filepath.Join(config.StoragePath, name)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of file %s from local storage: %w", name, err)
	}

	return content, nil
}
