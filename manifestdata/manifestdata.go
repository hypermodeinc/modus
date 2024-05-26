/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifestdata

import (
	"context"

	"hmruntime/logger"
	"hmruntime/storage"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
)

const manifestFileName = "hypermode.json"

var Manifest manifest.HypermodeManifest = manifest.HypermodeManifest{}

func MonitorManifestFile(ctx context.Context) {
	loadFile := func(file storage.FileInfo) error {
		if file.Name != manifestFileName {
			return nil
		}
		err := loadManifest(ctx)
		if err == nil {
			logger.Info(ctx).
				Str("filename", file.Name).
				Msg("Loaded manifest file.")
		} else {
			logger.Err(ctx, err).
				Str("filename", file.Name).
				Msg("Failed to load manifest file.")
		}

		return err
	}

	// NOTE: Removing the manifest file entirely is not currently supported.

	sm := storage.NewStorageMonitor(".json")
	sm.Added = loadFile
	sm.Modified = loadFile
	sm.Start(ctx)
}

func loadManifest(ctx context.Context) error {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	bytes, err := storage.GetFileContents(ctx, manifestFileName)
	if err != nil {
		return err
	}

	man, err := manifest.ReadManifest(bytes)
	if err != nil {
		return err
	}

	if !man.IsCurrentVersion() {
		logger.Warn(ctx).
			Str("filename", manifestFileName).
			Int("manifest_version", man.Version).
			Msg("The manifest file is in a deprecated format.  Please update it to the current format.")
	}

	// Only update the Manifest global when we have successfully read the manifest.
	Manifest = man
	return nil
}
