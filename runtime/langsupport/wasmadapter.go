/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package langsupport

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

type WasmAdapter interface {
	TypeInfo() LanguageTypeInfo
	Memory() wasm.Memory
	AllocateMemory(ctx context.Context, size uint32) (uint32, utils.Cleaner, error)
	GetFunction(name string) wasm.Function
	PreInvoke(ctx context.Context, plan ExecutionPlan) error
}
