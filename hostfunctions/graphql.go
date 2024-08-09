/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/graphqlclient"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func init() {
	addHostFunction("executeGQL", hostExecuteGQL, withI32Params(3), withI32Result())
}

func hostExecuteGQL(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName, stmt, varsJson string
	if err := readParams(ctx, mod, stack, &hostName, &stmt, &varsJson); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Executing GraphQL operation.",
		Completed: "Completed GraphQL operation.",
		Cancelled: "Cancelled GraphQL operation.",
		Error:     "Error executing GraphQL operation.",
		Detail:    fmt.Sprintf("Host: %s Query: %s", hostName, stmt),
	}

	// Prepare the host function
	var result string
	fn := func() (err error) {
		result, err = graphqlclient.Execute(ctx, hostName, stmt, varsJson)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, result); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}
