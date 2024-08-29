/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"context"

	"github.com/jensneuse/abstractlogger"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
)

type HypDSFactory struct {
	Ctx context.Context
}

func (f *HypDSFactory) Planner(logger abstractlogger.Logger) plan.DataSourcePlanner[HypDSConfig] {
	return &HypDSPlanner{
		ctx: f.Ctx,
	}
}

func (f *HypDSFactory) Context() context.Context {
	return f.Ctx
}
