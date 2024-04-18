/*
 * Copyright 2024 Hypermode, Inc.
 */

package engine

import (
	"context"
	"sync"
	"unsafe"

	"github.com/jensneuse/abstractlogger"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/plan"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

// Workaround for https://github.com/wundergraph/graphql-go-tools/issues/775
// Note is very unsafe and only works because there is no usage of the resolver
// in NewExecutionEngine before it's returned.

func newExecutionEngine(ctx context.Context, logger abstractlogger.Logger, engineConfig engine.Configuration,
	resolverOptions resolve.ResolverOptions) (*engine.ExecutionEngine, error) {

	// Create a new execution engine as usual
	e, err := engine.NewExecutionEngine(ctx, logger, engineConfig)
	if err != nil {
		return nil, err
	}

	// Create a new resolver so we can set different options than the default
	r := resolve.New(ctx, resolverOptions)

	// depends on the internal structure of engine.ExecutionEngine
	// the offset is calculated by the size of the fields before the resolver pointer
	offset := unsafe.Sizeof(logger)
	offset += unsafe.Sizeof(engineConfig)
	offset += unsafe.Sizeof(&plan.Planner{})
	offset += unsafe.Sizeof(sync.Mutex{})

	// find the resolver field and replace it with the new resolver
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(e)) + offset)
	*(*unsafe.Pointer)(ptr) = unsafe.Pointer(r)

	return e, nil
}
