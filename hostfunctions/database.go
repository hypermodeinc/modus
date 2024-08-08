/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/logger"
	"hmruntime/sqlclient"

	wasm "github.com/tetratelabs/wazero/api"
)

func init() {
	addHostFunction("databaseQuery", hostDatabaseQuery, withI32Params(4), withI32Result())
}

func hostDatabaseQuery(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName, dbType, statement, paramsJson string
	if err := readParams(ctx, mod, stack, &hostName, &dbType, &statement, &paramsJson); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Starting database query.",
		Completed: "Completed database query.",
		Cancelled: "Cancelled database query.",
		Error:     "Error querying database.",
		Detail:    fmt.Sprintf("Host: %s Query: %s", hostName, statement),
	}

	// Prepare the host function
	var response *sqlclient.HostQueryResponse
	fn := func() (err error) {
		response, err = sqlclient.ExecuteQuery(ctx, hostName, dbType, statement, paramsJson)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, response); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}
