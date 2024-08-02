/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/sqlclient"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostDatabaseQuery(ctx context.Context, mod wasm.Module, pHostName, pDbType, pStatement, pParamsJson uint32) uint32 {
	var hostName, dbType, statement, paramsJson string
	if err := readParams(ctx, mod, param{pHostName, &hostName}, param{pDbType, &dbType}, param{pStatement, &statement}, param{pParamsJson, &paramsJson}); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	response, err := sqlclient.ExecuteQuery(ctx, hostName, dbType, statement, paramsJson)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing database query.")
		return 0
	}

	offset, err := writeResult(ctx, mod, *response)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
