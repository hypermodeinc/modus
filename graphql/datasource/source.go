/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"hmruntime/functions"
	"hmruntime/host"
	"hmruntime/logger"
	"hmruntime/utils"

	"github.com/rs/xid"
)

const DataSourceName = "HypermodeFunctionsDataSource"

type Source struct{}

var errCallingFunction = fmt.Errorf("error calling function")

func (s Source) Load(ctx context.Context, input []byte, writer io.Writer) error {
	err := s.load(ctx, input, writer)
	if err != nil && !errors.Is(err, errCallingFunction) {
		// note: function call errors are already logged, so we don't log them again here
		logger.Err(ctx, err).Msg("Failed to load data.")
	}
	return err
}

func (s Source) load(ctx context.Context, input []byte, writer io.Writer) error {
	type callInfo struct {
		Function   string         `json:"fn"`
		Alias      string         `json:"alias"`
		Parameters map[string]any `json:"data"`
	}

	// Get the call info
	var ci callInfo
	dec := json.NewDecoder(bytes.NewReader(input))
	dec.UseNumber()
	err := dec.Decode(&ci)

	if err != nil {
		return fmt.Errorf("error getting function input: %w", err)
	}

	// Get the function info
	info, ok := functions.Functions[ci.Function]
	if !ok {
		return fmt.Errorf("no function registered named %s", ci.Function)
	}

	// Add plugin to the context
	ctx = context.WithValue(ctx, utils.PluginContextKey, info.Plugin)

	// Add execution ID to the context
	executionId := xid.New().String()
	ctx = context.WithValue(ctx, utils.ExecutionIdContextKey, executionId)

	// TODO: We should return the execution id(s) in the response with X-Hypermode-ExecutionID headers.
	// There might be multiple execution ids if the request triggers multiple function calls.

	// Get a module instance for this request.
	// Each request will get its own instance of the plugin module,
	// so that we can run multiple requests in parallel without risk
	// of corrupting the module's memory.
	mod, _, err := host.GetModuleInstance(ctx, info.Plugin)
	if err != nil {
		return fmt.Errorf("error getting module instance: %w", err)
	}
	defer mod.Close(ctx)

	// Call the function
	result, err := functions.CallFunction(ctx, mod, info, ci.Parameters)
	if err != nil {
		// TODO: pass buf somewhere
		return fmt.Errorf("%w '%s': %w", errCallingFunction, ci.Function, err)
		// return fmt.Errorf("error calling function '%s': %w", ci.Function, err)
	}

	// Write the result
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, `{"%s":%s}`, ci.Alias, jsonResult)

	return nil
}
