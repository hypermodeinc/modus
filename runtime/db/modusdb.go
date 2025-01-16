/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modusdb"
)

var dataDir = "data"
var GlobalModusDbEngine *modusdb.Engine

func InitModusDb(ctx context.Context) {
	// Initialize the database connection.
	var err error
	GlobalModusDbEngine, err = modusdb.NewEngine(modusdb.NewDefaultConfig(dataDir))
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to initialize modusdb.")
		return
	}
}

func CloseModusDb(ctx context.Context) {
	if GlobalModusDbEngine != nil {
		GlobalModusDbEngine.Close()
	}
}
