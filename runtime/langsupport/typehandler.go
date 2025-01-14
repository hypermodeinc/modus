/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package langsupport

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/utils"
)

type TypeHandler interface {
	TypeInfo() TypeInfo
	Read(ctx context.Context, wa WasmAdapter, offset uint32) (any, error)
	Write(ctx context.Context, wa WasmAdapter, offset uint32, obj any) (utils.Cleaner, error)
	Decode(ctx context.Context, wa WasmAdapter, vals []uint64) (any, error)
	Encode(ctx context.Context, wa WasmAdapter, obj any) ([]uint64, utils.Cleaner, error)
}
