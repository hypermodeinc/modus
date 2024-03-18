/*
 * Copyright 2024 Hypermode, Inc.
 */

package appdata

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/logger"
	"hmruntime/storage"
)

func LoadAppDataFiles(ctx context.Context) error {
	files, err := storage.ListFiles(ctx, ".json")
	if err != nil {
		return fmt.Errorf("failed to list application data files: %w", err)
	}

	for _, file := range files {
		err := loadAppData(ctx, file.Name)
		if err != nil {
			logger.Err(ctx, err).
				Str("filename", file.Name).
				Msg("Failed to load application data file.")
		}
	}

	return nil
}

func loadAppData(ctx context.Context, filename string) error {
	bytes, err := storage.GetFileContents(ctx, filename)
	if err != nil {
		return err
	}

	_, ok := appDataFiles[filename]
	if ok {
		err = json.Unmarshal(bytes, appDataFiles[filename])
		if err != nil {
			return err
		}
	}

	return nil
}
