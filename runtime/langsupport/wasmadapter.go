/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
