/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

import (
	"context"

	"hypruntime/plugins/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

type Planner interface {
	GetPlan(ctx context.Context, fnMeta *metadata.Function, fnDef wasm.FunctionDefinition) (ExecutionPlan, error)
	GetHandler(ctx context.Context, typ string) (TypeHandler, error)
	AllHandlers() map[string]TypeHandler
}
