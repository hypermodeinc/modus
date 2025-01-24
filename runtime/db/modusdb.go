/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"context"
	"path/filepath"
	"runtime"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modusdb"
)

var GlobalModusDbEngine *modusdb.Engine

func InitModusDb(ctx context.Context) {
	if config.IsDevEnvironment() && runtime.GOOS != "windows" {
		dataDir := filepath.Join(app.ModusHomeDir(), "data")
		if eng, err := modusdb.NewEngine(modusdb.NewDefaultConfig(dataDir)); err != nil {
			logger.Fatal(ctx).Err(err).Msg("Failed to initialize modusdb.")
		} else {
			GlobalModusDbEngine = eng
		}
	}
}

func CloseModusDb(ctx context.Context) {
	if GlobalModusDbEngine != nil {
		GlobalModusDbEngine.Close()
	}
}
