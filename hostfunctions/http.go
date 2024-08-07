/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/httpclient"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostFetch(ctx context.Context, mod wasm.Module, pRequest uint32) (pResponse uint32) {
	var request httpclient.HttpRequest
	if err := readParams(ctx, mod, param{pRequest, &request}); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	response, err := httpclient.HttpFetch(ctx, request)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error performing HTTP request.")
		return 0
	}

	offset, err := writeResult(ctx, mod, *response)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}
