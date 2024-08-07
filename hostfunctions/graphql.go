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

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pHostName, pStmt, pVarsJson uint32) uint32 {

	// Read input parameters
	var hostName, stmt, varsJson string
	err := readParams(ctx, mod, param{pHostName, &hostName}, param{pStmt, &stmt}, param{pVarsJson, &varsJson})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
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
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
