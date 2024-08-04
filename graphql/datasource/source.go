/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"hmruntime/functions"
	"hmruntime/logger"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	"github.com/buger/jsonparser"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/datasource/httpclient"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

const DataSourceName = "HypermodeFunctionsDataSource"

type callInfo struct {
	Function   fieldInfo      `json:"fn"`
	Parameters map[string]any `json:"data"`
}

type Source struct{}

func (s Source) Load(ctx context.Context, input []byte, out *bytes.Buffer) error {

	// Parse the input to get the function call info
	var ci callInfo
	err := utils.JsonDeserialize(input, &ci)
	if err != nil {
		return fmt.Errorf("error parsing input: %w", err)
	}

	// Load the data
	result, gqlErrors, err := s.callFunction(ctx, ci)

	// Write the response
	err = writeGraphQLResponse(out, result, gqlErrors, err, ci)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Error creating GraphQL response.")
	}

	return err
}

func (Source) LoadWithFiles(ctx context.Context, input []byte, files []httpclient.File, out *bytes.Buffer) (err error) {
	// See https://github.com/wundergraph/graphql-go-tools/pull/758
	panic("not implemented")
}

func (s Source) callFunction(ctx context.Context, callInfo callInfo) (any, []resolve.GraphQLError, error) {
	// Call the function
	info, err := wasmhost.CallFunctionWithParametersMap(ctx, callInfo.Function.Name, callInfo.Parameters)
	if err != nil {
		// The full error message has already been logged.  Return a generic error to the caller, which will be included in the response.
		return nil, nil, errors.New("error calling function")
	}

	// Store the execution info into the function output map.
	outputMap := ctx.Value(utils.FunctionOutputContextKey).(map[string]*wasmhost.ExecutionInfo)
	outputMap[callInfo.Function.AliasOrName()] = info

	// Transform messages (and error lines in the output buffers) to GraphQL errors.
	messages := append(info.Messages, utils.TransformConsoleOutput(info.Buffers)...)
	gqlErrors := transformErrors(messages, callInfo)

	return info.Result, gqlErrors, err
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

	// If there are GraphQL errors, serialize them as json
	var jsonErrors []byte
	if len(gqlErrors) > 0 {
		var err error
		jsonErrors, err = utils.JsonSerialize(gqlErrors)
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
	jsonResult, err := utils.JsonSerialize(result)
	if err != nil {
		return err
	}

	// Transform the data
	jsonData, err := transformData(jsonResult, &ci.Function)
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

func transformData(data []byte, tf *fieldInfo) ([]byte, error) {
	val, err := transformValue(data, tf)
	if err != nil {
		return nil, err
	}

	out := []byte(`{}`)
	return jsonparser.Set(out, val, tf.AliasOrName())
}

var nullWord = []byte("null")

func transformValue(data []byte, tf *fieldInfo) (result []byte, err error) {
	if len(tf.Fields) == 0 || len(data) == 0 || bytes.Equal(data, nullWord) {
		return data, nil
	}

	if tf.IsMapType {
		return transformMap(data, tf)
	}

	switch data[0] {
	case '{':
		return transformObject(data, tf)
	case '[':
		return transformArray(data, tf)
	default:
		return nil, fmt.Errorf("expected object or array")
	}
}

func transformArray(data []byte, tf *fieldInfo) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('[')

	var loopErr error
	_, err := jsonparser.ArrayEach(data, func(val []byte, _ jsonparser.ValueType, _ int, _ error) {
		if loopErr != nil {
			return
		}
		val, err := transformValue(val, tf)
		if err != nil {
			loopErr = err
			return
		}
		if buf.Len() > 1 {
			buf.WriteByte(',')
		}
		buf.Write(val)
	})
	if err != nil {
		return nil, err
	}
	if loopErr != nil {
		return nil, loopErr
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func transformObject(data []byte, tf *fieldInfo) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('{')
	for i, f := range tf.Fields {
		var val []byte
		if f.Name == "__typename" {
			val = []byte(`"` + tf.TypeName + `"`)
		} else {
			v, dataType, _, err := jsonparser.Get(data, f.Name)
			if err != nil {
				return nil, err
			}
			if dataType == jsonparser.String {
				// Note, string values here will be escaped for internal quotes, newlines, etc.,
				// but will be missing outer quotes.  So we need to add them back.
				v = []byte(`"` + string(v) + `"`)
			}
			val, err = transformValue(v, f)
			if err != nil {
				return nil, err
			}
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
	return buf.Bytes(), nil
}

func transformMap(data []byte, tf *fieldInfo) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('[')
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if buf.Len() > 1 {
			buf.WriteByte(',')
		}

		b := bytes.Buffer{}
		b.WriteByte('{')
		b.WriteString(`"key":"`)
		b.Write(key)
		b.WriteString(`","value":`)
		if dataType == jsonparser.String {
			b.WriteString(`"`)
			b.Write(value)
			b.WriteString(`"`)
		} else {
			b.Write(value)
		}
		b.WriteByte('}')

		val, err := transformObject(b.Bytes(), tf)
		if err != nil {
			return err
		}
		buf.Write(val)

		return nil
	})

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func transformErrors(messages []utils.LogMessage, ci callInfo) []resolve.GraphQLError {
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
