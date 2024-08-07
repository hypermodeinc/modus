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

func init() {
	addHostFunction(&hostFunctionDefinition{
		name:     "httpFetch",
		function: wasm.GoModuleFunc(hostHttpFetch),
		params:   []wasm.ValueType{wasm.ValueTypeI32},
		results:  []wasm.ValueType{wasm.ValueTypeI32},
	})
}

func hostHttpFetch(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var request httpclient.HttpRequest
	if err := readParams(ctx, mod, stack, &request); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, response); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}
