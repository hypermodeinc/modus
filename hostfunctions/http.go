/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/httpclient"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostFetch(ctx context.Context, mod wasm.Module, pRequest uint32) uint32 {

	// Read input parameters
	var request httpclient.HttpRequest
	if err := readParams(ctx, mod, param{pRequest, &request}); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Starting HTTP request.",
		Completed: "Completed HTTP request.",
		Cancelled: "Cancelled HTTP request.",
		Error:     "Error making HTTP request.",
		Detail:    fmt.Sprintf("%s %s", request.Method, request.Url),
	}

	// Prepare the host function
	var response *httpclient.HttpResponse
	fn := func() (err error) {
		response, err = httpclient.HttpFetch(ctx, request)
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
