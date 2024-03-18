/*
 * Copyright 2023 Hypermode, Inc.
 */

package host

import (
	"context"
	"fmt"

	"hmruntime/logger"
	"hmruntime/storage"
)

func LoadJsons(ctx context.Context) error {
	_, err := loadJsons(ctx)
	return err
}

func loadJsons(ctx context.Context) (map[string]bool, error) {
	var loaded = make(map[string]bool)

	files, err := storage.ListFiles(ctx, ".json")
	if err != nil {
		return nil, fmt.Errorf("failed to list JSON files: %w", err)
	}

	for _, file := range files {
		err := loadJson(ctx, file.Name)
		if err != nil {
			logger.Err(ctx, err).
				Str("filename", file.Name).
				Msg("Failed to load JSON file.")
		} else {
			loaded[file.Name] = true
		}
	}

	return loaded, nil
}
