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

	wasm "github.com/tetratelabs/wazero/api"
)

func hostFetchGet(ctx context.Context, mod wasm.Module, pHostName uint32, pStmt uint32) uint32 {

	var hostName, stmt string
	err := readParams2(ctx, mod, pHostName, pStmt, &hostName, &stmt)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	host, err := hosts.GetHost(hostName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting host.")
		return 0
	}

	result, err := connections.FetchGet[string](ctx, host, stmt)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing FetchGet.")
		return 0
	}

	offset, err := assemblyscript.WriteString(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}
	return offset
}
