/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifest

import (
	"context"

	"hmruntime/logger"
	"hmruntime/storage"
	"hmruntime/utils"

	"github.com/tailscale/hujson"
)

func MonitorAppDataFiles() {
	loadFile := func(file storage.FileInfo) error {
		ctx := context.Background()
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
	sm.Start()
}

func loadAppData(ctx context.Context, filename string) error {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	bytes, err := storage.GetFileContents(ctx, filename)
	if err != nil {
		return err
	}

	_, ok := manifestFiles[filename]
	if ok {

		// We allow comments and trailing commas in the JSON files.
		// This removes them, resulting in standard JSON.
		bytes, err := standardizeJSON(bytes)
		if err != nil {
			return err
		}

		// Now parse the JSON
		err = utils.JsonDeserialize(bytes, manifestFiles[filename])
		if err != nil {
			return err
		}
	}

	return nil
}

func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}
