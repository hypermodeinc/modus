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

func hostDatabaseQuery(ctx context.Context, mod wasm.Module, pHostName, pDbType, pStatement, pParamsJson uint32) uint32 {

	// Read input parameters
	var hostName, dbType, statement, paramsJson string
	if err := readParams(ctx, mod, param{pHostName, &hostName}, param{pDbType, &dbType}, param{pStatement, &statement}, param{pParamsJson, &paramsJson}); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
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
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, *response)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
