/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/graphqlclient"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pHostName uint32, pStmt uint32, pVarsJson uint32) uint32 {

	var hostName, stmt, varsJson string
	err := readParams(ctx, mod, param{pHostName, &hostName}, param{pStmt, &stmt}, param{pVarsJson, &varsJson})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	result, err := graphqlclient.Execute(ctx, hostName, stmt, varsJson)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error executing GraphQL operation.")
		return 0
	}

	offset, err := writeResult(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
