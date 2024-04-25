/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"hmruntime/functions"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	"github.com/buger/jsonparser"
	"github.com/rs/xid"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

const DataSourceName = "HypermodeFunctionsDataSource"

type callInfo struct {
	Function   templateField  `json:"fn"`
	Parameters map[string]any `json:"data"`
}

type FunctionOutput struct {
	ExecutionId string
	Buffers     utils.OutputBuffers
}

type Source struct{}

func (s Source) Load(ctx context.Context, input []byte, writer io.Writer) error {

	// Parse the input to get the function call info
	callInfo, err := parseInput(input)
	if err != nil {
		return fmt.Errorf("error parsing input: %w", err)
	}

	// Load the data
	result, gqlErrors, err := s.callFunction(ctx, callInfo)

	// Write the response
	return writeGraphQLResponse(writer, result, gqlErrors, err, callInfo)
}

func (s Source) callFunction(ctx context.Context, callInfo callInfo) (any, []resolve.GraphQLError, error) {

	// Get the function info
	info, ok := functions.Functions[callInfo.Function.Name]
	if !ok {
		return nil, nil, fmt.Errorf("no function registered named %s", callInfo.Function)
	}

	// Prepare the context that will be used throughout the function execution
	executionId := xid.New().String()
	ctx = context.WithValue(ctx, utils.ExecutionIdContextKey, executionId)
	ctx = context.WithValue(ctx, utils.PluginContextKey, info.Plugin)

	// Create output buffers for the function to write logs to
	buffers := utils.OutputBuffers{}

	// Get a module instance for this request.
	// Each request will get its own instance of the plugin module, so that we can run
	// multiple requests in parallel without risk of corrupting the module's memory.
	// This also protects against security risk, as each request will have its own
	// isolated memory space.  (One request cannot access another request's memory.)
	mod, err := wasmhost.GetModuleInstance(ctx, info.Plugin, &buffers)
	if err != nil {
		return nil, nil, err
	}
	defer mod.Close(ctx)

	// Call the function
	result, err := functions.CallFunction(ctx, mod, info, callInfo.Parameters)

	// Store the Execution ID and output buffers in the context
	outputMap := ctx.Value(utils.FunctionOutputContextKey).(map[string]FunctionOutput)
	outputMap[callInfo.Function.AliasOrName()] = FunctionOutput{
		ExecutionId: executionId,
		Buffers:     buffers,
	}

	// Transform error lines in the output buffers to GraphQL errors
	gqlErrors := transformErrors(buffers, callInfo)

	return result, gqlErrors, err
}

func parseInput(input []byte) (callInfo, error) {
	dec := json.NewDecoder(bytes.NewReader(input))
	dec.UseNumber()

	var ci callInfo
	err := dec.Decode(&ci)
	return ci, err
}

func writeGraphQLResponse(writer io.Writer, result any, gqlErrors []resolve.GraphQLError, fnErr error, ci callInfo) error {

	// Include the function error (except any we've filtered out)
	if fnErr != nil && functions.ShouldReturnErrorToResponse(fnErr) {
		gqlErrors = append(gqlErrors, resolve.GraphQLError{
			Message: fnErr.Error(),
			Path:    []any{ci.Function.AliasOrName()},
			Extensions: map[string]interface{}{
				"level": "error",
			},
		})
	}

	// If there are GraphQL errors, marshal them to json
	var jsonErrors []byte
	if len(gqlErrors) > 0 {
		var err error
		jsonErrors, err = json.Marshal(gqlErrors)
		if err != nil {
			return err
		}

		// If there are no other results, return only the errors
		if result == nil {
			fmt.Fprintf(writer, `{"errors":%s}`, jsonErrors)
			return nil
		}
	}

	// Get the data as json from the result
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// Transform the data
	jsonData, err = transformData(jsonData, ci.Function)
	if err != nil {
		return err
	}

	// Build and write the response, including errors if there are any
	if len(gqlErrors) > 0 {
		fmt.Fprintf(writer, `{"data":%s,"errors":%s}`, jsonData, jsonErrors)
	} else {
		fmt.Fprintf(writer, `{"data":%s}`, jsonData)
	}

	return nil
}

func transformData(data []byte, tf templateField) ([]byte, error) {
	val, err := transformValue(data, tf)
	if err != nil {
		return nil, err
	}

	out := []byte(`{}`)
	return jsonparser.Set(out, val, tf.AliasOrName())
}

func transformValue(data []byte, tf templateField) ([]byte, error) {
	if len(tf.Fields) == 0 || len(data) == 0 {
		return data, nil
	}

	buf := bytes.Buffer{}

	switch data[0] {
	case '{': // object
		buf.WriteByte('{')
		for i, f := range tf.Fields {
			val, dataType, _, err := jsonparser.Get(data, f.Name)
			if err != nil {
				return nil, err
			}
			if dataType == jsonparser.String {
				val, err = json.Marshal(string(val))
				if err != nil {
					return nil, err
				}
			}
			val, err = transformValue(val, f)
			if err != nil {
				return nil, err
			}
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('"')
			buf.WriteString(f.AliasOrName())
			buf.WriteString(`":`)
			buf.Write(val)
		}
		buf.WriteByte('}')

	case '[': // array
		buf.WriteByte('[')
		_, err := jsonparser.ArrayEach(data, func(val []byte, _ jsonparser.ValueType, _ int, _ error) {
			if buf.Len() > 1 {
				buf.WriteByte(',')
			}
			val, err := transformValue(val, tf)
			if err != nil {
				return
			}
			buf.Write(val)
		})
		if err != nil {
			return nil, err
		}

		buf.WriteByte(']')

	default:
		return nil, fmt.Errorf("expected object or array")
	}

	return buf.Bytes(), nil
}

func transformErrors(buffers utils.OutputBuffers, ci callInfo) []resolve.GraphQLError {
	messages := utils.TransformConsoleOutput(buffers)
	errors := make([]resolve.GraphQLError, 0, len(messages))
	for _, msg := range messages {
		// Only include errors.  Other messages will be captured later and
		// passed back as logs in the extensions section of the response.
		if msg.IsError() {
			errors = append(errors, resolve.GraphQLError{
				Message: msg.Message,
				Path:    []any{ci.Function.AliasOrName()},
				Extensions: map[string]interface{}{
					"level": msg.Level,
				},
			})
		}
	}
	return errors
}
