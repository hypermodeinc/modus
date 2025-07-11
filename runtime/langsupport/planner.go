/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package langsupport

import (
	"context"

	"github.com/hypermodeinc/modus/lib/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

type Planner interface {
	GetPlan(ctx context.Context, fnMeta *metadata.Function, fnDef wasm.FunctionDefinition) (ExecutionPlan, error)
	GetHandler(ctx context.Context, typ string) (TypeHandler, error)
	AllHandlers() map[string]TypeHandler
}
