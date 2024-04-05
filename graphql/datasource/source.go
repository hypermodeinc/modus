/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"hmruntime/functions"
	"hmruntime/host"
	"hmruntime/logger"

	"github.com/rs/xid"
)

const DataSourceName = "HypermodeFunctionsDataSource"

type Source struct{}

func (s Source) Load(ctx context.Context, input []byte, writer io.Writer) error {
	type callInfo struct {
		Function   string         `json:"fn"`
		Alias      string         `json:"alias"`
		Parameters map[string]any `json:"data"`
	}

	// Get the call info
	var ci callInfo
	err := json.Unmarshal(input, &ci)
	if err != nil {
		return fmt.Errorf("error getting function input: %w", err)
	}

	// Get the function info
	// TODO: Are all functions query resolvers? what about mutations?
	resolver := "Query." + ci.Function
	info, ok := functions.FunctionsMap[resolver]
	if !ok {
		return fmt.Errorf("no function registered for %s", resolver)
	}

	// Add plugin details to the context
	ctx = context.WithValue(ctx, logger.PluginNameContextKey, info.Plugin.Name())
	ctx = context.WithValue(ctx, logger.BuildIdContextKey, info.Plugin.BuildId())

	// Add execution ID to the context
	executionId := xid.New().String()
	ctx = context.WithValue(ctx, logger.ExecutionIdContextKey, executionId)

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
		return fmt.Errorf("error calling function '%s': %w", ci.Function, err)
	}

	// Determine if the result is already JSON
	isJson := false
	fieldType := info.Schema.FieldDef.Type.NamedType
	if _, ok := result.(string); ok && fieldType != "String" {
		isJson = true
	}

	// Write the result
	if isJson {
		fmt.Fprintf(writer, `{"%s":%s}`, ci.Alias, result)
	} else {
		jsonResult, err := json.Marshal(result)
		if err != nil {
			return err
		}
		fmt.Fprintf(writer, `{"%s":%s}`, ci.Alias, jsonResult)
	}

	return nil
}
