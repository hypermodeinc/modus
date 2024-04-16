/*
 * Copyright 2024 Hypermode, Inc.
 */

package appdata

import (
	"context"
	"encoding/json"

	"hmruntime/logger"
	"hmruntime/storage"
)

func MonitorAppDataFiles(ctx context.Context) {

	loadFile := func(file storage.FileInfo) error {
		err := loadAppData(ctx, file.Name)
		if err == nil {
			logger.Info(ctx).
				Str("filename", file.Name).
				Msg("Loaded application data file.")
		} else {
			logger.Err(ctx, err).
				Str("filename", file.Name).
				Msg("Failed to load application data file.")
		}

		return err
	}

	// NOTE: Removing a file entirely is not currently supported.

	sm := storage.NewStorageMonitor(".json")
	sm.Added = loadFile
	sm.Modified = loadFile
	sm.Start(ctx)
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
