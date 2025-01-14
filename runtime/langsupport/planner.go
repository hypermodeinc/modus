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

	"github.com/hypermodeinc/modus/lib/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

type Planner interface {
	GetPlan(ctx context.Context, fnMeta *metadata.Function, fnDef wasm.FunctionDefinition) (ExecutionPlan, error)
	GetHandler(ctx context.Context, typ string) (TypeHandler, error)
	AllHandlers() map[string]TypeHandler
}
