/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package envfiles

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/storage"
)

func MonitorEnvFiles(ctx context.Context) {
	loadFile := func(file storage.FileInfo) error {
		logger.Info(ctx).Str("filename", file.Name).Msg("Env file changed. Loading environment variables.")
		err := LoadEnvFiles(ctx)
		if err != nil {
			logger.Err(ctx, err).Str("filename", file.Name).Msg("Failed to load env file.")
			return err
		}

		return nil
	}
	sm := storage.NewStorageMonitor(".env", ".env.*")
	sm.Added = loadFile
	sm.Changed = func(errors []error) {
		logger.Info(ctx).Msg("Env files changed. Updating environment variables.")
		if len(errors) == 0 {
			err := LoadEnvFiles(ctx)
			if err != nil {
				logger.Err(ctx, err).Msg("Failed to load env files.")
			}
		}
	}

	sm.Start(ctx)
}
