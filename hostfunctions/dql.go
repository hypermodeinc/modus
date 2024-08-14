/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/dqlclient"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func init() {
	addHostFunction("executeDQL", hostExecuteDQL, withI32Params(3), withI32Result())
}

func hostExecuteDQL(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName, stmt, varsJson string
	if err := readParams(ctx, mod, stack, &hostName, &stmt, &varsJson); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Executing DQL operation.",
		Completed: "Completed DQL operation.",
		Cancelled: "Cancelled DQL operation.",
		Error:     "Error executing DQL operation.",
		Detail:    fmt.Sprintf("Host: %s Query: %s", hostName, stmt),
	}

	// Prepare the host function
	var result string
	fn := func() (err error) {
		result, err = dqlclient.ExecuteQuery(ctx, hostName, stmt, varsJson)
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
