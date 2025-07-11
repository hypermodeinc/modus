/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package datasource

import (
	"context"

	"github.com/jensneuse/abstractlogger"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/ast"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
)

func NewModusDataSourceFactory(ctx context.Context) plan.PlannerFactory[ModusDataSourceConfig] {
	return &modusDataSourceFactory{
		ctx: ctx,
	}
}

type modusDataSourceFactory struct {
	ctx context.Context
}

func (f *modusDataSourceFactory) Planner(logger abstractlogger.Logger) plan.DataSourcePlanner[ModusDataSourceConfig] {
	return &modusDataSourcePlanner{
		ctx: f.ctx,
	}
}

func (f *modusDataSourceFactory) Context() context.Context {
	return f.ctx
}

func (f *modusDataSourceFactory) UpstreamSchema(dataSourceConfig plan.DataSourceConfiguration[ModusDataSourceConfig]) (*ast.Document, bool) {
	return nil, false
}
