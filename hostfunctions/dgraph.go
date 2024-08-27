/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/dgraphclient"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func init() {
	addHostFunction("executeDQL", hostExecuteDQL, withI32Params(2), withI32Result())
	addHostFunction("dgraphAlterSchema", hostDgraphAlterSchema, withI32Params(2), withI32Result())
	addHostFunction("dgraphDropAttr", hostDgraphDropAttr, withI32Params(2), withI32Result())
	addHostFunction("dgraphDropAll", hostDgraphDropAll, withI32Params(1), withI32Result())
}

func hostExecuteDQL(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName string
	var req dgraphclient.Request
	if err := readParams(ctx, mod, stack, &hostName, &req); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Executing DQL operation.",
		Completed: "Completed DQL operation.",
		Cancelled: "Cancelled DQL operation.",
		Error:     "Error executing DQL operation.",
		Detail:    fmt.Sprintf("Host: %s Req: %s", hostName, fmt.Sprint(req)),
	}

	// Prepare the host function
	var result *dgraphclient.Response
	fn := func() (err error) {
		result, err = dgraphclient.Execute(ctx, hostName, &req)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, result); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}

}

func hostDgraphAlterSchema(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName, schema string
	if err := readParams(ctx, mod, stack, &hostName, &schema); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Altering DQL schema.",
		Completed: "Completed DQL schema alteration.",
		Cancelled: "Cancelled DQL schema alteration.",
		Error:     "Error altering DQL schema.",
		Detail:    fmt.Sprintf("Host: %s Schema: %s", hostName, schema),
	}

	// Prepare the host function
	var result string
	fn := func() (err error) {
		result, err = dgraphclient.AlterSchema(ctx, hostName, schema)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, result); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostDgraphDropAttr(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName, attr string
	if err := readParams(ctx, mod, stack, &hostName, &attr); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Dropping DQL attribute.",
		Completed: "Completed DQL attribute drop.",
		Cancelled: "Cancelled DQL attribute drop.",
		Error:     "Error dropping DQL attribute.",
		Detail:    fmt.Sprintf("Host: %s Attribute: %s", hostName, attr),
	}

	// Prepare the host function
	var result string
	fn := func() (err error) {
		result, err = dgraphclient.DropAttr(ctx, hostName, attr)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, result); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostDgraphDropAll(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName string
	if err := readParams(ctx, mod, stack, &hostName); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Dropping all DQL data.",
		Completed: "Completed DQL data drop.",
		Cancelled: "Cancelled DQL data drop.",
		Error:     "Error dropping DQL data.",
		Detail:    fmt.Sprintf("Host: %s", hostName),
	}

	// Prepare the host function
	var result string
	fn := func() (err error) {
		result, err = dgraphclient.DropAll(ctx, hostName)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, result); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}
