/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package datasource

import (
	"github.com/hypermodeinc/modus/runtime/wasmhost"
)

type HypDSConfig struct {
	WasmHost          wasmhost.WasmHost
	FieldsToFunctions map[string]string
	MapTypes          []string
}
