/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package envfiles

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/storage"
)

func MonitorEnvFiles(ctx context.Context) {
	sm := storage.NewStorageMonitor(".env", ".env.*")
	sm.Changed = func(errors []error) {
		logger.Info(ctx).Msg("Env files changed. Updating environment variables.")
		if len(errors) == 0 {
			err := LoadEnvFiles(ctx)
			if err != nil {
				const msg = "Failed to load env files."
				sentryutils.CaptureError(ctx, err, msg)
				logger.Error(ctx, err).Msg(msg)
			}
		}
	}

	sm.Start(ctx)
}
