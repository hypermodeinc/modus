/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"encoding/json"

	"hmruntime/connections"
	"hmruntime/functions/assemblyscript"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pHostName uint32, pStmt uint32, pVars uint32) uint32 {
	mem := mod.Memory()

	host, err := getHost(mem, pHostName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting host.")
		return 0
	}

	stmt, err := assemblyscript.ReadString(mem, pStmt)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading GraphQL query string from wasm memory.")
		return 0
	}

	sVars, err := assemblyscript.ReadString(mem, pVars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading GraphQL variables string from wasm memory.")
		return 0
	}

	vars := make(map[string]any)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		logger.Err(ctx, err).Msg("Error unmarshalling GraphQL variables.")
		return 0
	}

	result, err := connections.ExecuteGraphqlApi[string](ctx, host, stmt, vars)
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
