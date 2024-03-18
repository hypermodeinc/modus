/*
 * Copyright 2024 Hypermode, Inc.
 */

package storage

import (
	"context"
	"hmruntime/config"
	"hmruntime/logger"
	"time"
)

type StorageMonitor struct {
	extension string
	files     map[string]*monitoredFile
	Added     func(FileInfo)
	Modified  func(FileInfo)
	Removed   func(FileInfo)
	Changed   func()
}

type monitoredFile struct {
	file     FileInfo
	lastSeen time.Time
}

func NewStorageMonitor(extension string) *StorageMonitor {
	return &StorageMonitor{
		extension: extension,
		files:     make(map[string]*monitoredFile),
		Added:     func(FileInfo) {},
		Modified:  func(FileInfo) {},
		Removed:   func(FileInfo) {},
		Changed:   func() {},
	}
}

func (sm *StorageMonitor) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(config.RefreshInterval)
		defer ticker.Stop()

		var loggedError = false

		for {
			files, err := impl.listFiles(ctx, sm.extension)
			if err != nil {
				// Don't stop watching. We'll just try again on the next cycle.
				if !loggedError {
					logger.Err(ctx, err).Msgf("Failed to list %s files.", sm.extension)
					loggedError = true
				}
				continue
			} else {
				loggedError = false
			}

			// Compare list of files retrieved to existing files
			var changed = false
			var thisTime = time.Now()
			for _, file := range files {
				existing, found := sm.files[file.Name]
				if !found {
					// New file
					changed = true
					sm.files[file.Name] = &monitoredFile{file, thisTime}
					sm.Added(file)
				} else if file.Hash != existing.file.Hash ||
					(file.Hash == "" && file.LastModified.After(existing.file.LastModified)) {
					// Modified file
					changed = true
					sm.files[file.Name] = &monitoredFile{file, thisTime}
					sm.Modified(file)
				} else {
					// No change
					existing.lastSeen = thisTime
				}
			}

			// Check for removed files
			for name, file := range sm.files {
				if file.lastSeen.Before(thisTime) {
					changed = true
					delete(sm.files, name)
					sm.Removed(file.file)
				}
			}

			// Notify if anything changed
			if changed {
				sm.Changed()
			}

			// Handle auto-reload disabled
			if config.NoReload {
				return
			}

			// Wait for next cycle
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return
			}
		}
	}()
}
