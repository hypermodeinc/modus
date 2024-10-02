/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/plugins/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

type Planner interface {
	GetPlan(ctx context.Context, fnMeta *metadata.Function, fnDef wasm.FunctionDefinition) (ExecutionPlan, error)
	GetHandler(ctx context.Context, typ string) (TypeHandler, error)
	AllHandlers() map[string]TypeHandler
}
