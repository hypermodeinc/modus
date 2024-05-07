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

func hostFetch(ctx context.Context, mod wasm.Module, pHostName uint32, pMethod uint32, pPath uint32, pBody uint32, pHeaders uint32) uint32 {

	var hostName, method, path, body string
	var headers map[string]string
	err := readParams5(ctx, mod, pHostName, pMethod, pPath, pBody, pHeaders, &hostName, &method, &path, &body, &headers)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	host, err := hosts.GetHost(hostName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting host.")
		return 0
	}

	result, err := connections.Fetch[string](ctx, host, method, path, body, headers)
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
