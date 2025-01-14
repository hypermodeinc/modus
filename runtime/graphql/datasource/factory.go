/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package datasource

import (
	"context"

	"github.com/jensneuse/abstractlogger"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/ast"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
)

func NewHypDSFactory(ctx context.Context) plan.PlannerFactory[HypDSConfig] {
	return &hypDSFactory{
		ctx: ctx,
	}
}

type hypDSFactory struct {
	ctx context.Context
}

func (f *hypDSFactory) Planner(logger abstractlogger.Logger) plan.DataSourcePlanner[HypDSConfig] {
	return &HypDSPlanner{
		ctx: f.ctx,
	}
}

func (f *hypDSFactory) Context() context.Context {
	return f.ctx
}

func (f *hypDSFactory) UpstreamSchema(dataSourceConfig plan.DataSourceConfiguration[HypDSConfig]) (*ast.Document, bool) {
	return nil, false
}
