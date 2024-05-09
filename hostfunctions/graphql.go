/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/connections"
	"hmruntime/functions/assemblyscript"
	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pHostName uint32, pStmt uint32, pVars uint32) uint32 {

	var hostName, stmt, sVars string
	err := readParams3(ctx, mod, pHostName, pStmt, pVars, &hostName, &stmt, &sVars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	vars := make(map[string]any)
	if err := utils.JsonDeserialize([]byte(sVars), &vars); err != nil {
		logger.Err(ctx, err).Msg("Error deserializing GraphQL variables.")
		return 0
	}

	host, err := hosts.GetHost(hostName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting host.")
		return 0
	}

	result, err := connections.ExecuteGraphqlApi(ctx, host, stmt, vars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing GraphQL operation.")
		return 0
	}

	offset, err := assemblyscript.WriteString(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
