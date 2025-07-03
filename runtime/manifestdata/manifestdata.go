/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifestdata

import (
	"context"
	"sync"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/storage"
	"github.com/hypermodeinc/modus/runtime/utils"
)

const manifestFileName = "modus.json"

var mu sync.RWMutex
var man = &manifest.Manifest{}

func GetManifest() *manifest.Manifest {
	mu.RLock()
	defer mu.RUnlock()
	return man
}

func SetManifest(m *manifest.Manifest) {
	mu.Lock()
	defer mu.Unlock()
	man = m
}

func MonitorManifestFile(ctx context.Context) {
	loadFile := func(file storage.FileInfo) error {
		logger.Info(ctx).Str("filename", file.Name).Msg("Loading manifest file.")
		if file.Name != manifestFileName {
			return nil
		}

		if err := loadManifest(ctx); err != nil {
			logger.Error(ctx, err).Str("filename", file.Name).Msg("Failed to load manifest file.")
			return err
		}

		return nil
	}

	sm := storage.NewStorageMonitor("*.json")
	sm.Added = loadFile
	sm.Modified = loadFile
	sm.Removed = func(file storage.FileInfo) error {
		if file.Name == manifestFileName {
			logger.Warn(ctx).Str("filename", file.Name).Msg("Manifest file removed.")
			if err := unloadManifest(ctx); err != nil {
				logger.Error(ctx, err).Str("filename", file.Name).Msg("Failed to unload manifest file.")
				return err
			}
		}
		return nil
	}
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
	SetManifest(m)

	return triggerManifestLoaded(ctx)
}

func unloadManifest(ctx context.Context) error {
	SetManifest(&manifest.Manifest{})
	return triggerManifestLoaded(ctx)
}
