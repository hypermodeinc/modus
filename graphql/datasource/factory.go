/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"context"

	"github.com/jensneuse/abstractlogger"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
)

type Factory[T Configuration] struct {
	Ctx context.Context
}

func (f *Factory[T]) Planner(logger abstractlogger.Logger) plan.DataSourcePlanner[T] {
	return &Planner[T]{
		ctx: f.Ctx,
	}
}

func (f *Factory[T]) Context() context.Context {
	return f.Ctx
}
