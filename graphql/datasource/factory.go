/*
 * Copyright 2024 Hypermode, Inc.
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
