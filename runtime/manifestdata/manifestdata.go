/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifestdata

import (
	"context"
	"sync"

	"hypruntime/logger"
	"hypruntime/storage"
	"hypruntime/utils"

	"github.com/hypermodeAI/manifest"
)

const manifestFileName = "hypermode.json"

var mu sync.RWMutex
var man = &manifest.HypermodeManifest{}

func GetManifest() *manifest.HypermodeManifest {
	mu.RLock()
	defer mu.RUnlock()
	return man
}

func SetManifest(m *manifest.HypermodeManifest) {
	mu.Lock()
	defer mu.Unlock()
	man = m
}

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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	bytes, err := storage.GetFileContents(ctx, manifestFileName)
	if err != nil {
		return err
	}

	m, err := manifest.ReadManifest(bytes)
	if err != nil {
		return err
	}

	if !m.IsCurrentVersion() {
		logger.Warn(ctx).
			Str("filename", manifestFileName).
			Int("manifest_version", m.Version).
			Msg("The manifest file is in a deprecated format.  Please update it to the current format.")
	}

	// Only update the Manifest global when we have successfully read the manifest.
	SetManifest(&m)

	// Trigger the manifest loaded event.
	err = triggerManifestLoaded(ctx)

	return err
}
