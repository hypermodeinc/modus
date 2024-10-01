/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

import (
	"context"

	"hypruntime/utils"
)

type TypeHandler interface {
	TypeInfo() TypeInfo
	Read(ctx context.Context, wa WasmAdapter, offset uint32) (any, error)
	Write(ctx context.Context, wa WasmAdapter, offset uint32, obj any) (utils.Cleaner, error)
	Decode(ctx context.Context, wa WasmAdapter, vals []uint64) (any, error)
	Encode(ctx context.Context, wa WasmAdapter, obj any) ([]uint64, utils.Cleaner, error)
}
