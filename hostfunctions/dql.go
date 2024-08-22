/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/dqlclient"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func init() {
	addHostFunction("executeDQLQuery", hostExecuteDQLQuery, withI32Params(3), withI32Result())
	addHostFunction("executeDQLMutations", hostExecuteDQLMutations, withI32Params(3), withI32Result())
	addHostFunction("executeDQLUpserts", hostExecuteDQLUpserts, withI32Params(4), withI32Result())
	addHostFunction("dgraphAlterSchema", hostDgraphAlterSchema, withI32Params(2), withI32Result())
	addHostFunction("dgraphDropAttr", hostDgraphDropAttr, withI32Params(2), withI32Result())
	addHostFunction("dgraphDropAll", hostDgraphDropAll, withI32Params(1), withI32Result())
}

func hostExecuteDQLQuery(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName, query, varsJson string
	var mutations []string
	if err := readParams(ctx, mod, stack, &hostName, &query, &varsJson); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Executing DQL operation.",
		Completed: "Completed DQL operation.",
		Cancelled: "Cancelled DQL operation.",
		Error:     "Error executing DQL operation.",
		Detail:    fmt.Sprintf("Host: %s Query: %s Mutations: %v", hostName, query, mutations),
	}

	// Prepare the host function
	var result string
	fn := func() (err error) {
		result, err = dqlclient.ExecuteQuery(ctx, hostName, query, varsJson)
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

func hostExecuteDQLMutations(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName string
	var setMutations, delMutations []string
	if err := readParams(ctx, mod, stack, &hostName, &setMutations, &delMutations); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Executing DQL mutations.",
		Completed: "Completed DQL mutations.",
		Cancelled: "Cancelled DQL mutations.",
		Error:     "Error executing DQL mutations.",
		Detail:    fmt.Sprintf("Host: %s SetMutations: %v DelMutations: %v", hostName, setMutations, delMutations),
	}

	// Prepare the host function
	var result map[string]string
	fn := func() (err error) {
		result, err = dqlclient.ExecuteMutations(ctx, hostName, setMutations, delMutations)
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

func hostExecuteDQLUpserts(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var hostName, query string
	var setMutations, delMutations []string
	if err := readParams(ctx, mod, stack, &hostName, &query, &setMutations, &delMutations); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Executing DQL upserts.",
		Completed: "Completed DQL upserts.",
		Cancelled: "Cancelled DQL upserts.",
		Error:     "Error executing DQL upserts.",
		Detail:    fmt.Sprintf("Host: %s Query: %s SetMutations: %v DelMutations: %v", hostName, query, setMutations, delMutations),
	}

	// Prepare the host function
	var result map[string]string
	fn := func() (err error) {
		result, err = dqlclient.ExecuteUpserts(ctx, hostName, query, setMutations, delMutations)
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
		result, err = dqlclient.AlterSchema(ctx, hostName, schema)
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
		result, err = dqlclient.DropAttr(ctx, hostName, attr)
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
		result, err = dqlclient.DropAll(ctx, hostName)
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
